package main

import (
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	gctx "github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/ijalalfrz/majoo-test-1/auth"
	"github.com/ijalalfrz/majoo-test-1/config"
	"github.com/ijalalfrz/majoo-test-1/jwt"
	"github.com/ijalalfrz/majoo-test-1/middleware"
	"github.com/ijalalfrz/majoo-test-1/response"
	"github.com/ijalalfrz/majoo-test-1/server"
	"github.com/ijalalfrz/majoo-test-1/transaction"
	_ "github.com/joho/godotenv/autoload" // for development
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmlogrus"
)

var (
	cfg          *config.Config
	location     *time.Location
	indexMessage string = "Application is running properly"
)

func init() {
	cfg = config.Load()
}

func main() {
	// set logger
	logger := logrus.New()
	logger.SetFormatter(cfg.Logger.Formatter)
	logger.SetReportCaller(true)
	logger.AddHook(&apmlogrus.Hook{
		LogLevels: logrus.AllLevels,
	})

	// set validator
	vld := validator.New()

	// set mariadb object
	db, err := sql.Open(cfg.Mariadb.Driver, cfg.Mariadb.DSN)
	if err != nil {
		logger.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		logger.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(cfg.Mariadb.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.Mariadb.MaxIdleConnections)

	// set jwt
	secret := []byte("majoo-test-secret")
	jwtUser := jwt.NewJSONWebToken(secret)

	// middleware
	sessionMiddleware := middleware.NewSessionMiddleware(jwtUser)

	// set router object
	router := mux.NewRouter()
	router.HandleFunc("/majoo-service", index)

	// set domain object
	transactionRepository := transaction.NewTransactionRepository(logger, db)
	transactionUsecase := transaction.NewTransactionUsecase(transaction.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		Logger:      logger,
		Repository:  transactionRepository,
	})
	transaction.NewTransactionHTTPHandler(logger, sessionMiddleware, router, transactionUsecase)

	authRepository := auth.NewAuthRepository(logger, db)
	authUsecase := auth.NewAuthUsecase(auth.UsecaseProperty{
		ServiceName: cfg.Application.Name,
		Logger:      logger,
		Repository:  authRepository,
		JWT:         jwtUser,
	})
	auth.NewAuthHTTPHandler(logger, vld, router, authUsecase)

	// middleware
	httpHandler := gctx.ClearHandler(router)
	httpHandler = middleware.Recovery(logger, httpHandler)
	httpHandler = middleware.CORS(httpHandler)

	// initiate server
	srv := server.NewServer(logger, httpHandler, cfg.Application.Port)
	srv.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt)
	<-sigterm

	// closing service for a gracefull shutdown.
	db.Close()
	srv.Close()
}

func index(w http.ResponseWriter, r *http.Request) {
	resp := response.NewSuccessResponse(nil, response.StatOK, indexMessage)
	response.JSON(w, resp)
}
