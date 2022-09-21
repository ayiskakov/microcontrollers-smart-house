package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"microcontrollers/internal/handler"
	"microcontrollers/internal/server"
	"microcontrollers/internal/service"
	"microcontrollers/internal/storage/memstorage"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	storage := memstorage.NewStorage()
	svc := service.NewService(storage)
	hnd := handler.NewHandler(svc)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), hnd.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Print("Home Automation Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Home Automation Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
