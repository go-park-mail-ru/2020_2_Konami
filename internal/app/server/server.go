package server

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	meetingDeliveryPkg "konami_backend/internal/pkg/meeting/delivery/http"
	meetingRepoPkg "konami_backend/internal/pkg/meeting/repository"
	meetingUseCasePkg "konami_backend/internal/pkg/meeting/usecase"
	profileDeliveryPkg "konami_backend/internal/pkg/profile/delivery/http"
	profileRepoPkg "konami_backend/internal/pkg/profile/repository"
	profileUseCasePkg "konami_backend/internal/pkg/profile/usecase"
	sessionDeliveryPkg "konami_backend/internal/pkg/session/delivery/http"
	sessionRepoPkg "konami_backend/internal/pkg/session/repository"
	sessionUseCasePkg "konami_backend/internal/pkg/session/usecase"
	tagRepoPkg "konami_backend/internal/pkg/tag/repository"
	uploadsHandlerPkg "konami_backend/internal/pkg/utils/uploads_handler"
	"net/http"
	"os"
)

func InitDelivery(db *gorm.DB, uploadsDir, meetPicsDir, userPicsDir, defMeetPic, defUserPic string) (
	meetingDeliveryPkg.MeetingHandler,
	profileDeliveryPkg.ProfileHandler,
	sessionDeliveryPkg.SessionHandler,
) {
	meetingRepo := meetingRepoPkg.NewMeetingGormRepo(db)
	profileRepo := profileRepoPkg.NewProfileGormRepo(db)
	sessionRepo := sessionRepoPkg.NewSessionGormRepo(db)
	tagRepo := tagRepoPkg.NewTagGormRepo(db)
	uploadsHandler := uploadsHandlerPkg.NewUploadsHandler(uploadsDir)
	meetingUC := meetingUseCasePkg.NewMeetingUseCase(
		meetingRepo, uploadsHandler, tagRepo, meetPicsDir, defMeetPic)
	profileUC := profileUseCasePkg.NewProfileUseCase(
		profileRepo, uploadsHandler, tagRepo, userPicsDir, defUserPic)
	sessionUC := sessionUseCasePkg.NewSessionUseCase(sessionRepo)
	meetingDelivery := meetingDeliveryPkg.MeetingHandler{
		MeetingUC: meetingUC,
		SessionUC: sessionUC,
	}
	profileDelivery := profileDeliveryPkg.ProfileHandler{
		ProfileUC: profileUC,
		SessionUC: sessionUC,
	}
	sessionDelivery := sessionDeliveryPkg.SessionHandler{
		SessionUC: sessionUC,
		ProfileUC: profileUC,
	}
	return meetingDelivery, profileDelivery, sessionDelivery
}

func InitRouter(db *gorm.DB, uploadsDir, meetPicsDir, userPicsDir, defMeetPic, defUserPic string) http.Handler {
	meeting, profile, session := InitDelivery(db, uploadsDir, meetPicsDir, userPicsDir, defMeetPic, defUserPic)
	r := mux.NewRouter()
	r.HandleFunc("/api/meetings", meeting.GetMeetingsList).Methods("GET")
	r.HandleFunc("/api/meeting", meeting.CreateMeeting).Methods("POST")
	r.HandleFunc("/api/meet", meeting.GetMeeting).Methods("GET")
	r.HandleFunc("/api/meet", meeting.UpdateMeeting).Methods("POST")

	r.HandleFunc("/api/people", profile.GetPeople).Methods("GET")
	r.HandleFunc("/api/user", profile.GetUser).Methods("GET")
	r.HandleFunc("/api/user", profile.EditUser).Methods("POST")
	r.HandleFunc("/api/images", profile.UploadUserPic).Methods("POST")
	r.HandleFunc("/api/signup", profile.SignUp).Methods("POST")

	r.HandleFunc("/api/me", session.GetUserId).Methods("GET")
	r.HandleFunc("/api/login", session.LogIn).Methods("POST")
	r.HandleFunc("/api/logout", session.LogOut).Methods("POST")
	return r
}

func Start() {
	db, err := gorm.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatalf("Failed to launch db: %v", err)
	}
	defer db.Close()
	if err := db.DB().Ping(); err != nil {
		log.Fatalf("Failed to launch db: %v", err)
	}

	r := InitRouter(db, "uploads", "meetingpics", "userpics",
		"assets/paris.jpg", "assets/empty-avatar.jpeg")

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
		err = http.ListenAndServe(":"+port, r)
	} else {
		log.Println("Launching at HTTPS port " + tlsPort)
		err = http.ListenAndServeTLS(":"+tlsPort, certFile, keyFile, r)
	}

	if err != nil {
		log.Fatal("Unable to launch server: ", err)
	}
}
