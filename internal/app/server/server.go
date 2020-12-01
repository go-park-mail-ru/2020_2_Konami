package server

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	meetingDeliveryPkg "konami_backend/internal/pkg/meeting/delivery/http"
	meetingRepoPkg "konami_backend/internal/pkg/meeting/repository"
	meetingUseCasePkg "konami_backend/internal/pkg/meeting/usecase"
	messageDeliveryPkg "konami_backend/internal/pkg/message/delivery/http"
	messageRepoPkg "konami_backend/internal/pkg/message/repository"
	messageUseCasePkg "konami_backend/internal/pkg/message/usecase"
	"konami_backend/internal/pkg/middleware"
	profileDeliveryPkg "konami_backend/internal/pkg/profile/delivery/http"
	profileRepoPkg "konami_backend/internal/pkg/profile/repository"
	profileUseCasePkg "konami_backend/internal/pkg/profile/usecase"
	tagRepoPkg "konami_backend/internal/pkg/tag/repository"
	corsInit "konami_backend/internal/pkg/utils/cors_init"
	"konami_backend/internal/pkg/utils/token_handler"
	uploadsHandlerPkg "konami_backend/internal/pkg/utils/uploads_handler"
	loggerPkg "konami_backend/logger"
	authProto "konami_backend/proto/auth"
	csrfProto "konami_backend/proto/csrf"
	"log"
	"net/http"
	"os"
	"strconv"
)

func InitDelivery(db *gorm.DB, log *loggerPkg.Logger, maxReqSize int64,
	authClient authProto.AuthCheckerClient,
	csrfClient csrfProto.CsrfDispatcherClient,
	uploadsDir, meetPicsDir, userPicsDir, defMeetPic, defUserPic string) (
	meetingDeliveryPkg.MeetingHandler,
	profileDeliveryPkg.ProfileHandler,
	messageDeliveryPkg.MessageHandler,
	token_handler.TokenHandler,
	middleware.AuthMiddleware,
	middleware.CSRFMiddleware,
	middleware.AccessLogMiddleware,
	error,
) {
	profileRepo := profileRepoPkg.NewProfileGormRepo(db)
	meetingRepo := meetingRepoPkg.NewMeetingGormRepo(db, profileRepo)
	tagRepo := tagRepoPkg.NewTagGormRepo(db)
	msgRepo := messageRepoPkg.NewMeetingGormRepo(db)
	uploadsHandler := uploadsHandlerPkg.NewUploadsHandler(uploadsDir)
	meetingUC := meetingUseCasePkg.NewMeetingUseCase(
		meetingRepo, uploadsHandler, tagRepo, meetPicsDir, defMeetPic)
	profileUC := profileUseCasePkg.NewProfileUseCase(
		profileRepo, uploadsHandler, tagRepo, userPicsDir, defUserPic)
	msgUC := messageUseCasePkg.NewMessageUseCase(msgRepo)
	meetingDelivery := meetingDeliveryPkg.MeetingHandler{
		MeetingUC:  meetingUC,
		MaxReqSize: maxReqSize,
	}
	profileDelivery := profileDeliveryPkg.ProfileHandler{
		ProfileUC:  profileUC,
		AuthClient: authClient,
		MaxReqSize: maxReqSize,
	}
	tokenHandler := token_handler.TokenHandler{CsrfClient: csrfClient, Log: log}
	msgDelivery := messageDeliveryPkg.NewMessageHandler(msgUC, log, maxReqSize)
	authM := middleware.NewAuthMiddleware(profileUC, authClient)
	csrfM := middleware.NewCsrfMiddleware(csrfClient, log)
	logM := middleware.NewAccessLogMiddleware(log)
	return meetingDelivery, profileDelivery, msgDelivery, tokenHandler, authM, csrfM, logM, nil
}

func InitRouter(
	meeting meetingDeliveryPkg.MeetingHandler,
	profile profileDeliveryPkg.ProfileHandler,
	message messageDeliveryPkg.MessageHandler,
	token token_handler.TokenHandler,
	authM middleware.AuthMiddleware,
	csrfM middleware.CSRFMiddleware,
	logM middleware.AccessLogMiddleware,
	panicM middleware.PanicMiddleware) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", message.Upgrade)
	rApi := mux.NewRouter()
	r.PathPrefix("/api/").Handler(http.StripPrefix("/api", rApi))
	rApi.HandleFunc("/people", profile.GetPeople).Methods("GET")
	rApi.HandleFunc("/user", profile.GetUser).Methods("GET")
	rApi.HandleFunc("/signup", profile.SignUp).Methods("POST")
	rApi.HandleFunc("/login", profile.LogIn).Methods("POST")
	rApi.HandleFunc("/csrf", token.GetCSRF).Methods("GET")

	rApi.HandleFunc("/meeting", meeting.GetMeeting).Methods("GET")
	rApi.HandleFunc("/meetings", meeting.GetMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/my", meeting.GetUserMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/favorite", meeting.GetFavMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/top", meeting.GetTopMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/recommended", meeting.GetRecommendedList).Methods("GET")
	rApi.HandleFunc("/meetings/tagged", meeting.GetTaggedMeetings).Methods("GET")
	rApi.HandleFunc("/meetings/akin", meeting.GetAkinMeetings).Methods("GET")
	rApi.HandleFunc("/meetings/search", meeting.SearchMeetings).Methods("GET")

	rApi.HandleFunc("/me", profile.GetUserId).Methods("GET")
	rApi.HandleFunc("/logout", profile.LogOut).Methods("DELETE")
	rApi.HandleFunc("/meeting", meeting.CreateMeeting).Methods("POST")
	rApi.HandleFunc("/meeting", meeting.UpdateMeeting).Methods("PATCH")
	rApi.HandleFunc("/user", profile.EditUser).Methods("PATCH")
	rApi.HandleFunc("/images", profile.UploadUserPic).Methods("POST")

	r.Handle("/metrics", promhttp.Handler())

	rApi.HandleFunc("/messages", message.GetMessages).Methods("GET")
	rApi.HandleFunc("/message", message.SendMessage).Methods("POST")
	go message.ServeWS()

	r.Use(panicM.PanicRecovery)
	r.Use(middleware.HeadersMiddleware)
	rApi.Use(logM.Log)
	rApi.Use(authM.Auth)
	rApi.Use(csrfM.CSRFCheck)
	return r
}

func Start() {
	var logger *loggerPkg.Logger
	logger = loggerPkg.NewLogger(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
	dsn := os.Getenv("DB_CONN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatalf("failed to launch db: %v", err)
	}
	dbdb, err := db.DB()
	if err != nil {
		logger.Fatalf("failed to launch db: %v", err)
	}
	defer dbdb.Close()
	if err := dbdb.Ping(); err != nil {
		logger.Fatalf("failed to launch db: %v", err)
	}

	authAddr := os.Getenv("AUTH_ADDR")
	if authAddr == "" {
		authAddr = "127.0.0.1:8002"
	}
	authConn, err := grpc.Dial(
		authAddr,
		grpc.WithUnaryInterceptor(loggerPkg.GetGRPCInterceptor(logger)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("unable to connect to grpc")
	}
	defer authConn.Close()
	authClient := authProto.NewAuthCheckerClient(authConn)

	csrfAddr := os.Getenv("CSRF_ADDR")
	if csrfAddr == "" {
		csrfAddr = "127.0.0.1:8003"
	}
	csrfConn, err := grpc.Dial(
		csrfAddr,
		grpc.WithUnaryInterceptor(loggerPkg.GetGRPCInterceptor(logger)),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("unable to connect to grpc")
	}
	defer csrfConn.Close()
	csrfClient := csrfProto.NewCsrfDispatcherClient(csrfConn)

	maxReqSize, err := strconv.ParseInt(os.Getenv("MAX_REQ_SIZE"), 10, 64)
	if err != nil || maxReqSize <= 0 {
		maxReqSize = 10 * 1024 * 1024
	}

	meeting, profile, msg, token, authM, csrfM, logM, err := InitDelivery(
		db, logger, maxReqSize, authClient, csrfClient,
		"uploads", "meetingpics", "userpics",
		"assets/paris.jpg", "assets/empty-avatar.jpeg")
	if err != nil {
		logger.Fatalf("failed to init delivery: %v", err)
		return
	}

	panicM := middleware.NewPanicMiddleware(logger)
	r := InitRouter(meeting, profile, msg, token, authM, csrfM, logM, panicM)
	c := corsInit.InitCors()
	h := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	certFile := os.Getenv("CERTFILE")
	if certFile == "" {
		certFile = "/etc/letsencrypt/live/okto.pw/fullchain.pem"
	}
	keyFile := os.Getenv("KEYFILE")
	if keyFile == "" {
		keyFile = "/etc/letsencrypt/live/okto.pw/privkey.pem"
	}
	tlsHost := os.Getenv("TLSHOST")
	if tlsHost == "" {
		tlsHost = "okto.pw"
	}
	tlsPort := os.Getenv("TLSPORT")

	if tlsPort == "" {
		logger.Println("Launching at HTTP port " + port)
		err = http.ListenAndServe(":"+port, h)
	} else {
		logger.Println("Launching at HTTPS port " + tlsPort)
		err = http.ListenAndServeTLS(":"+tlsPort, certFile, keyFile, h)
	}

	if err != nil {
		logger.Fatal("Unable to launch server: ", err)
	}
}

func Migrate() {
	dsn := os.Getenv("DB_CONN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	dbdb, err := db.DB()
	if err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	defer dbdb.Close()
	if err := dbdb.Ping(); err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	err = db.AutoMigrate(
		&tagRepoPkg.Tag{},
		&profileRepoPkg.SkillTag{},
		&profileRepoPkg.InterestTag{},
		&profileRepoPkg.Profile{},
		&meetingRepoPkg.Registration{},
		&meetingRepoPkg.Like{},
		&meetingRepoPkg.Meeting{},
		&messageRepoPkg.Message{},
	)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
	db.Exec(`
CREATE INDEX IF NOT EXISTS search_idx ON meetings USING gin((
setweight(to_tsvector('russian', title), 'A') || setweight(to_tsvector('english', title), 'A') ||
setweight(to_tsvector('russian', text), 'B') || setweight(to_tsvector('english', text), 'B') || 
setweight(to_tsvector('russian', city), 'C') || setweight(to_tsvector('english', city), 'C') ||
setweight(to_tsvector('russian', address), 'D') || setweight(to_tsvector('english', address), 'D')
	));`)
	var tags = []tagRepoPkg.Tag{
		{Name: "ИТ и интернет"}, {Name: "Языки программирования"}, {Name: "C++"},
		{Name: "Python"}, {Name: "JavaScript"}, {Name: "Golang"}, {Name: "Mail.ru"},
		{Name: "Yandex"}, {Name: "Бизнес"}, {Name: "Хобби"}, {Name: "Творчество"},
		{Name: "Кино"}, {Name: "Театры"}, {Name: "Вечеринки"}, {Name: "Еда"}, {Name: "Концерты"},
		{Name: "Спорт"}, {Name: "Красота"}, {Name: "Здоровье"}, {Name: "Наука"},
		{Name: "Выставки"}, {Name: "Искусство"}, {Name: "Культура"}, {Name: "Экскурсии"},
		{Name: "Путешествия"}, {Name: "Психология"}, {Name: "Образование"}, {Name: "Россия"}}

	for _, tag := range tags {
		res := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&tag)
		if res.Error != nil {
			log.Fatalf("failed to create tags: %v", err)
		}
	}
}

func Truncate() {
	dsn := os.Getenv("DB_CONN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	dbdb, err := db.DB()
	if err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	defer dbdb.Close()
	if err := dbdb.Ping(); err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	db.Exec("DELETE FROM profile_skill_tags")
	db.Exec("DELETE FROM profile_interest_tags")
	db.Exec("DELETE FROM profile_meeting_tags")
	db.Exec("DELETE FROM meeting_tags")
	db.Exec("DELETE FROM registrations")
	db.Exec("DELETE FROM likes")
	db.Exec("DELETE FROM meetings")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM messages")
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&profileRepoPkg.InterestTag{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&profileRepoPkg.SkillTag{})
}
