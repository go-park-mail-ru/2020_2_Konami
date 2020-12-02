package auth_server

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	sessionDeliveryPkg "konami_backend/auth/pkg/session/delivery/grpc"
	sessionRepoPkg "konami_backend/auth/pkg/session/repository"
	sessionUseCasePkg "konami_backend/auth/pkg/session/usecase"
	loggerPkg "konami_backend/logger"
	"konami_backend/proto/auth"
	"log"
	"net"
	"os"
)

func InitDelivery(db *gorm.DB) (sessionDeliveryPkg.SessionHandler, error) {
	sessionRepo := sessionRepoPkg.NewSessionGormRepo(db)
	sessionUC := sessionUseCasePkg.NewSessionUseCase(sessionRepo)
	sessionDelivery := sessionDeliveryPkg.NewSessionHandler(sessionUC)
	return sessionDelivery, nil
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
	sessionHandler, err := InitDelivery(db)
	if err != nil {
		logger.Fatalf("failed to init delivery: %v", err)
		return
	}

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = "8002"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalln("cant listen port", err)
	}

	server := grpc.NewServer()
	auth.RegisterAuthCheckerServer(server, &sessionHandler)

	fmt.Println("starting server at :", port)
	err = server.Serve(lis)
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
		&sessionRepoPkg.Session{},
	)
	if err != nil {
		log.Fatalf("failed to migrate db: %v", err)
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
	db.Exec("DELETE FROM sessions")
}
