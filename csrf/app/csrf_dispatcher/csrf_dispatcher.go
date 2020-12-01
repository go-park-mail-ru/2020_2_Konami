package csrf_dispatcher

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	csrfDeliveryPkg "konami_backend/csrf/pkg/csrf/delivery/grpc"
	csrfRepoPkg "konami_backend/csrf/pkg/csrf/repository"
	csrfUseCasePkg "konami_backend/csrf/pkg/csrf/usecase"
	loggerPkg "konami_backend/logger"
	"konami_backend/proto/csrf"
	"net"
	"os"
	"strconv"
)

func InitDelivery(rconn *redis.Pool, log *loggerPkg.Logger, csrfSecret string, csrfExpire int64) (
	csrfDeliveryPkg.CsrfHandler, error,
) {
	csrfRepo := csrfRepoPkg.NewRedisTokenManager(rconn)
	csrfUC, err := csrfUseCasePkg.NewCsrfUseCase(csrfSecret, csrfExpire, csrfRepo)
	if err != nil {
		log.Error(err)
		return csrfDeliveryPkg.CsrfHandler{}, err
	}

	csrfDelivery := csrfDeliveryPkg.NewCsrfHandler(csrfUC, log)
	return csrfDelivery, nil
}

func Start() {
	var logger *loggerPkg.Logger
	logger = loggerPkg.NewLogger(os.Stdout)
	logger.SetLevel(logrus.TraceLevel)
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
				logger.Fatalf("failed to launch redis pool: %v", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()

	csrfSecret := os.Getenv("CSRF_SECRET")
	if csrfSecret == "" {
		logger.Fatalf("csrf secret not provided")
	}

	csrfDuration, err := strconv.ParseInt(os.Getenv("CSRF_DURATION"), 10, 64)
	if err != nil || csrfDuration <= 0 {
		csrfDuration = 3600
	}
	csrfHandler, err := InitDelivery(redisConn, logger, os.Getenv("CSRF_SECRET"), csrfDuration)
	if err != nil {
		logger.Fatalf("failed to init delivery: %v", err)
		return
	}
	port := os.Getenv("CSRF_PORT")
	if port == "" {
		port = "8003"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Fatalln("unable to listen port ", port, err)
	}

	server := grpc.NewServer()
	csrf.RegisterCsrfDispatcherServer(server, &csrfHandler)

	fmt.Println("starting server at :", port)
	err = server.Serve(lis)
	if err != nil {
		logger.Fatal("Unable to launch server: ", err)
	}
}
