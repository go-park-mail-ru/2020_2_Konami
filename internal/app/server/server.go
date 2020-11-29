package server

import (
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	csrfDeliveryPkg "konami_backend/internal/pkg/csrf/delivery/http"
	csrfRepoPkg "konami_backend/internal/pkg/csrf/repository"
	csrfUseCasePkg "konami_backend/internal/pkg/csrf/usecase"
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
	sessionDeliveryPkg "konami_backend/internal/pkg/session/delivery/http"
	sessionRepoPkg "konami_backend/internal/pkg/session/repository"
	sessionUseCasePkg "konami_backend/internal/pkg/session/usecase"
	tagRepoPkg "konami_backend/internal/pkg/tag/repository"
	corsInit "konami_backend/internal/pkg/utils/cors_init"
	uploadsHandlerPkg "konami_backend/internal/pkg/utils/uploads_handler"
	"konami_backend/logger"
	"log"
	"net/http"
	"os"
)

func InitDelivery(db *gorm.DB, rconn *redis.Pool, log *logger.Logger, maxReqSize int64,
	csrfSecret string, csrfExpire int64,
	uploadsDir, meetPicsDir, userPicsDir, defMeetPic, defUserPic string) (

	csrfDeliveryPkg.CSRFHandler,
	meetingDeliveryPkg.MeetingHandler,
	profileDeliveryPkg.ProfileHandler,
	sessionDeliveryPkg.SessionHandler,
	messageDeliveryPkg.MessageHandler,
	middleware.AuthMiddleware,
	middleware.CSRFMiddleware,
	middleware.AccessLogMiddleware,
	error,
) {
	csrfRepo := csrfRepoPkg.NewRedisTokenManager(rconn)
	meetingRepo := meetingRepoPkg.NewMeetingGormRepo(db)
	profileRepo := profileRepoPkg.NewProfileGormRepo(db)
	sessionRepo := sessionRepoPkg.NewSessionGormRepo(db)
	tagRepo := tagRepoPkg.NewTagGormRepo(db)
	msgRepo := messageRepoPkg.NewMeetingGormRepo(db)
	uploadsHandler := uploadsHandlerPkg.NewUploadsHandler(uploadsDir)
	csrfUC, err := csrfUseCasePkg.NewCsrfUseCase(csrfSecret, csrfExpire, csrfRepo)
	if err != nil {
		log.Error(err)
		return csrfDeliveryPkg.CSRFHandler{}, meetingDeliveryPkg.MeetingHandler{},
			profileDeliveryPkg.ProfileHandler{}, sessionDeliveryPkg.SessionHandler{},
			messageDeliveryPkg.MessageHandler{}, middleware.AuthMiddleware{},
			middleware.CSRFMiddleware{}, middleware.AccessLogMiddleware{}, err
	}
	meetingUC := meetingUseCasePkg.NewMeetingUseCase(
		meetingRepo, uploadsHandler, tagRepo, meetPicsDir, defMeetPic)
	profileUC := profileUseCasePkg.NewProfileUseCase(
		profileRepo, uploadsHandler, tagRepo, userPicsDir, defUserPic)
	sessionUC := sessionUseCasePkg.NewSessionUseCase(sessionRepo)
	msgUC := messageUseCasePkg.NewMessageUseCase(msgRepo)
	csrfDelivery := csrfDeliveryPkg.CSRFHandler{
		CsrfUC: csrfUC,
		Log:    log,
	}
	meetingDelivery := meetingDeliveryPkg.MeetingHandler{
		MeetingUC:  meetingUC,
		SessionUC:  sessionUC,
		MaxReqSize: maxReqSize,
	}
	profileDelivery := profileDeliveryPkg.ProfileHandler{
		ProfileUC:  profileUC,
		SessionUC:  sessionUC,
		MaxReqSize: maxReqSize,
	}
	sessionDelivery := sessionDeliveryPkg.SessionHandler{
		SessionUC: sessionUC,
		ProfileUC: profileUC,
	}
	msgDelivery := messageDeliveryPkg.MessageHandler{
		MessageUC:  msgUC,
		Log:        log,
		MaxReqSize: maxReqSize,
	}
	authM := middleware.NewAuthMiddleware(profileUC, sessionUC)
	csrfM := middleware.NewCsrfMiddleware(csrfUC, log)
	logM := middleware.NewAccessLogMiddleware(log)
	return csrfDelivery, meetingDelivery, profileDelivery, sessionDelivery, msgDelivery, authM, csrfM, logM, nil
}

func InitRouter(
	csrf csrfDeliveryPkg.CSRFHandler,
	meeting meetingDeliveryPkg.MeetingHandler,
	profile profileDeliveryPkg.ProfileHandler,
	session sessionDeliveryPkg.SessionHandler,
	message messageDeliveryPkg.MessageHandler,
	authM middleware.AuthMiddleware,
	csrfM middleware.CSRFMiddleware,
	logM middleware.AccessLogMiddleware,
	panicM middleware.PanicMiddleware) http.Handler {

	r := mux.NewRouter()
	rApi := mux.NewRouter()
	r.PathPrefix("/api/").Handler(http.StripPrefix("/api", rApi))
	rApi.HandleFunc("/people", profile.GetPeople).Methods("GET")
	rApi.HandleFunc("/user", profile.GetUser).Methods("GET")
	rApi.HandleFunc("/signup", profile.SignUp).Methods("POST")
	rApi.HandleFunc("/login", session.LogIn).Methods("POST")
	rApi.HandleFunc("/csrf", csrf.GetCSRF).Methods("GET")
	rApi.HandleFunc("/meeting", meeting.GetMeeting).Methods("GET")

	rApi.HandleFunc("/meetings", meeting.GetMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/my", meeting.GetUserMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/favorite", meeting.GetFavMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/top", meeting.GetTopMeetingsList).Methods("GET")
	rApi.HandleFunc("/meetings/recommended", meeting.GetRecommendedList).Methods("GET")
	rApi.HandleFunc("/meetings/tagged", meeting.GetTaggedMeetings).Methods("GET")
	rApi.HandleFunc("/meetings/akin", meeting.GetAkinMeetings).Methods("GET")
	rApi.HandleFunc("/meetings/search", meeting.GetSearchMeetings).Methods("GET")

	rApi.HandleFunc("/me", session.GetUserId).Methods("GET")
	rApi.HandleFunc("/logout", session.LogOut).Methods("DELETE")
	rApi.HandleFunc("/meeting", meeting.CreateMeeting).Methods("POST")
	rApi.HandleFunc("/meeting", meeting.UpdateMeeting).Methods("PATCH")
	rApi.HandleFunc("/user", profile.EditUser).Methods("PATCH")
	rApi.HandleFunc("/images", profile.UploadUserPic).Methods("POST")

	rApi.HandleFunc("/messages", message.GetMessages).Methods("GET")
	rApi.HandleFunc("/message", message.SendMessage).Methods("GET")
	rApi.HandleFunc("/ws", message.Upgrade)
	go message.ServeWS()

	r.Use(panicM.PanicRecovery)
	r.Use(middleware.HeadersMiddleware)
	r.Use(logM.Log)
	r.Use(authM.Auth)
	r.Use(csrfM.CSRFCheck)

	return r
}

func Start() {
	var log *logger.Logger
	log = logger.NewLogger(os.Stdout)
	log.SetLevel(logrus.TraceLevel)

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
	redisAddr := os.Getenv("REDIS_CONN")
	if redisAddr == "" {
		redisAddr = "redis://user:@localhost:6379/0"
	}
	redisConn := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(redisAddr)
			if err != nil {
				log.Fatalf("failed to launch redis pool: %v", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()

	csrfSecret := os.Getenv("CSRF_SECRET")
	if csrfSecret == "" {
		log.Fatalf("csrf secret not provided")

	}

	var maxRecSize int64 = 10 * 1024 * 1024
	var csrfDuration int64 = 3600
	csrf, meeting, profile, session, msg, authM, csrfM, logM, err := InitDelivery(
		db, redisConn, log, maxRecSize,
		os.Getenv("CSRF_SECRET"), csrfDuration,
		"uploads", "meetingpics", "userpics",
		"assets/paris.jpg", "assets/empty-avatar.jpeg")
	if err != nil {
		log.Fatalf("failed to init delivery: %v", err)
		return
	}

	panicM := middleware.NewPanicMiddleware(log)
	r := InitRouter(csrf, meeting, profile, session, msg, authM, csrfM, logM, panicM)
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
		log.Println("Launching at HTTP port " + port)
		err = http.ListenAndServe(":"+port, h)
	} else {
		log.Println("Launching at HTTPS port " + tlsPort)
		err = http.ListenAndServeTLS(":"+tlsPort, certFile, keyFile, h)
	}

	if err != nil {
		log.Fatal("Unable to launch server: ", err)
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
		&sessionRepoPkg.Session{},
		&messageRepoPkg.Message{},
	)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
	}
	var tags = []tagRepoPkg.Tag{
		{Name: "ИТ и интернет"}, {Name: "Языки программирования"}, {Name: "C++"},
		{Name: "Python"}, {Name: "JavaScript"}, {Name: "Golang"}, {Name: "Mail.ru"},
		{Name: "Yandex"}, {Name: "Бизнес"}, {Name: "Хобби"}, {Name: "Творчество"},
		{Name: "Кино"}, {Name: "Театры"}, {Name: "Вечеринки"}, {Name: "Еда"}, {Name: "Концерты"},
		{Name: "Спорт"}, {Name: "Красота"}, {Name: "Здоровье"}, {Name: "Наука"},
		{Name: "Выставки"}, {Name: "Искусство"}, {Name: "Культура"}, {Name: "Экскурсии"},
		{Name: "Путешествия"}, {Name: "Психология"}, {Name: "Образование"}, {Name: "Россия"}}

	for _, tag := range tags {
		res := db.Create(&tag)
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
	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM profiles")
	db.Exec("DELETE FROM tags")
	db.Exec("DELETE FROM messages")
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&profileRepoPkg.InterestTag{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&profileRepoPkg.SkillTag{})
}
