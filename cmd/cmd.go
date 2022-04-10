package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/sea-auca/auc-auth/config"
	"github.com/sea-auca/auc-auth/db"
	"github.com/sea-auca/auc-auth/user/repo"
	"github.com/sea-auca/auc-auth/user/server"
	"github.com/sea-auca/auc-auth/user/service"
	"github.com/talkanbaev-artur/shutdown"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := shutdown.NewShutdown()
	logger := InitLogger(shutdown) // init logger
	config.Init((*zap.Logger)(logger))

	//SERVICES INITIALISATION STEP

	_, err := db.ConnectDatabase(shutdown)
	pdb := db.ConnectPGXDatabase(ctx)
	serv, router := NewServer()
	if err != nil {
		logger.Fatal("Could not connect to the database", zap.Error(err))
	}
	userRepo, verLinkRepo := repo.NewPgxUserRepo(pdb), repo.NewPgxVerificationLinkRepo(pdb)
	userService := service.NewService(userRepo, verLinkRepo)

	server.RegisterRoutes(userService, router)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go Listen(logger, serv)

	go func() {
		oscall := <-c
		logger.Info("Catched shutdown sygnal", zap.Any("Signal", oscall))
		cancel()
	}()

	logger.Info("Initialised authentication server")
	<-ctx.Done()
	errs := shutdown.Close() // initialise the shutdown process and print all errors
	if len(errs) != 0 {
		for _, v := range errs {
			logger.Info("Shutdown error", zap.Error(v))
		}
	}
}

func InitLogger(shutdown *shutdown.Shutdown) *zap.Logger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	zap.ReplaceGlobals(logger)
	shutdown.Add(logger.Sync)
	return logger
}

func NewServer() (http.Server, *mux.Router) {
	conf := config.Config().Server
	r := mux.NewRouter()
	serv := http.Server{
		Addr:        fmt.Sprintf("%s:%d", "0.0.0.0", conf.Port),
		ReadTimeout: time.Second * time.Duration(conf.ReadTimeout),
		Handler:     r,
	}
	return serv, r
}

func Listen(logger *zap.Logger, serv http.Server) {
	logger.Error("Fatal http server error", zap.Error(serv.ListenAndServe()))
}
