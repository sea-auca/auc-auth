package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sea/auth/config"
	"sea/auth/utils"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// read configuration
	_, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	// create global defer object
	shutdown := utils.NewShutdown()
	logger := InitLogger(shutdown) // init logger

	ctx, cancel := context.WithCancel(context.Background())

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

func InitLogger(shutdown *utils.Shutdown) *zap.SugaredLogger {
	log1, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	logger := log1.Sugar()

	shutdown.Add(logger.Sync)
	return logger
}
