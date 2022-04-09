package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sea-auca/auc-auth/config"
	"github.com/sea-auca/auc-auth/db"
	"github.com/sea-auca/auc-auth/user/repo"
	"github.com/sea-auca/auc-auth/user/server"
	"github.com/sea-auca/auc-auth/user/service"
	"github.com/talkanbaev-artur/shutdown"
	"go.uber.org/zap"
)

func main() {
	// read configuration
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// create global defer object
	shutdown := shutdown.NewShutdown()
	logger := InitLogger(shutdown) // init logger
	ctx, cancel := context.WithCancel(context.Background())

	//SERVICES INITIALISATION STEP

	rep, err := db.ConnectDatabase(conf, shutdown)

	if err != nil {
		logger.Fatalw("Could not connect to the database", "Error", err)
	}
	userRepo, verLinkRepo := repo.NewUserRepository(rep), repo.NewVerificationRepository(rep)
	userService := service.NewService(userRepo, email, verLinkRepo)

	go server.NewHTTP(userService)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		oscall := <-c
		logger.Infow("Catched shutdown sygnal", "Signal", oscall)
		cancel()
	}()

	logger.Infow("Initialised authentication server", "ConfigStatus", "OK", "ServiceStatus", "OK")
	<-ctx.Done()
	errs := shutdown.Close() // initialise the shutdown process and print all errors
	if len(errs) != 0 {
		for _, v := range errs {
			logger.Infow("Shutdown error", "Error msg", v)
		}
	}
}

func InitLogger(shutdown *shutdown.Shutdown) *zap.SugaredLogger {
	log1, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := log1.Sugar()
	zap.ReplaceGlobals(log1)
	shutdown.Add(logger.Sync)
	return logger
}
