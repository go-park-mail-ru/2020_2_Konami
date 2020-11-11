package server

import (
	"flag"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	csrfDeliveryPkg "konami_backend/internal/pkg/csrf/delivery/http"
	csrfRepoPkg "konami_backend/internal/pkg/csrf/repository"
	csrfUseCasePkg "konami_backend/internal/pkg/csrf/usecase"
	meetingDeliveryPkg "konami_backend/internal/pkg/meeting/delivery/http"
	meetingRepoPkg "konami_backend/internal/pkg/meeting/repository"
	meetingUseCasePkg "konami_backend/internal/pkg/meeting/usecase"
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
	uploadsHandler := uploadsHandlerPkg.NewUploadsHandler(uploadsDir)
	csrfUC, err := csrfUseCasePkg.NewCsrfUseCase(csrfSecret, csrfExpire, csrfRepo)
	if err != nil {
		log.Error(err)
		return csrfDeliveryPkg.CSRFHandler{}, meetingDeliveryPkg.MeetingHandler{},
			profileDeliveryPkg.ProfileHandler{}, sessionDeliveryPkg.SessionHandler{},
			middleware.AuthMiddleware{}, middleware.CSRFMiddleware{}, middleware.AccessLogMiddleware{}, err
	}
	meetingUC := meetingUseCasePkg.NewMeetingUseCase(
		meetingRepo, uploadsHandler, tagRepo, meetPicsDir, defMeetPic)
	profileUC := profileUseCasePkg.NewProfileUseCase(
		profileRepo, uploadsHandler, tagRepo, userPicsDir, defUserPic)
	sessionUC := sessionUseCasePkg.NewSessionUseCase(sessionRepo)
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
	authM := middleware.NewAuthMiddleware(profileUC, sessionUC)
	csrfM := middleware.NewCsrfMiddleware(csrfUC, log)
	logM := middleware.NewAccessLogMiddleware(log)
	return csrfDelivery, meetingDelivery, profileDelivery, sessionDelivery, authM, csrfM, logM, nil
}

func InitRouter(
	csrf csrfDeliveryPkg.CSRFHandler,
	meeting meetingDeliveryPkg.MeetingHandler,
	profile profileDeliveryPkg.ProfileHandler,
	session sessionDeliveryPkg.SessionHandler,
	authM middleware.AuthMiddleware,
	csrfM middleware.CSRFMiddleware,
	logM middleware.AccessLogMiddleware,
	panicM middleware.PanicMiddleware) http.Handler {

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/people", profile.GetPeople).Methods("GET")
	r.HandleFunc("/api/v1/user", profile.GetUser).Methods("GET")
	r.HandleFunc("/api/v1/signup", profile.SignUp).Methods("POST")
	r.HandleFunc("/api/v1/login", session.LogIn).Methods("POST")
	r.HandleFunc("/api/v1/csrf", csrf.GetCSRF).Methods("GET")
	r.HandleFunc("/api/v1/meeting", meeting.GetMeeting).Methods("GET")
	r.HandleFunc("/api/v1/meetings", meeting.GetMeetingsList).Methods("GET")
	r.HandleFunc("/api/v1/me", session.GetUserId).Methods("GET")
	r.HandleFunc("/api/v1/logout", session.LogOut).Methods("DELETE")
	r.HandleFunc("/api/v1/meeting", meeting.CreateMeeting).Methods("POST")
	r.HandleFunc("/api/v1/meeting", meeting.UpdateMeeting).Methods("PATCH")
	r.HandleFunc("/api/v1/user", profile.EditUser).Methods("PATCH")
	r.HandleFunc("/api/v1/images", profile.UploadUserPic).Methods("POST")
	r.Use(panicM.PanicRecovery)
	r.Use(middleware.HeadersMiddleware)
	r.Use(logM.Log)
	r.Use(authM.Auth)
	rCSRF := r.Headers("Csrf-Token").Subrouter()
	rCSRF.Use(csrfM.CSRFCheck)
	return r
}

func Start() {
	var log *logger.Logger
	log = logger.NewLogger(os.Stdout)
	log.SetLevel(logrus.TraceLevel)

	db, err := gorm.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}
	defer db.Close()
	if err := db.DB().Ping(); err != nil {
		log.Fatalf("failed to launch db: %v", err)
	}

	redisAddr := flag.String("addr", os.Getenv("REDIS_CONN"), "redis addr")
	redisConn := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(*redisAddr)
			if err != nil {
				log.Fatalf("failed to launch redis pool: %v", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()

	var maxRecSize int64 = 10 * 1024 * 1024
	var csrfDuration int64 = 3600
	csrf, meeting, profile, session, authM, csrfM, logM, err := InitDelivery(
		db, redisConn, log, maxRecSize,
		os.Getenv("CSRF_SECRET"), csrfDuration,
		"uploads", "meetingpics", "userpics",
		"assets/paris.jpg", "assets/empty-avatar.jpeg")
	if err != nil {
		log.Fatalf("failed to init delivery: %v", err)
		return
	}

	panicM := middleware.NewPanicMiddleware(log)
	r := InitRouter(csrf, meeting, profile, session, authM, csrfM, logM, panicM)
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
