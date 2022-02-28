package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/polygate/api/router"
	"github.com/quanxiang-cloud/polygate/pkg/config"
)

var (
	configPath = flag.String("config", "../configs/config.yml", "-config config file path")
)

func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		logger.Logger.Panicw("", err)
	}

	logger.Logger = logger.New(conf.Log)

	router, err := router.NewRouter(conf)
	if err != nil {
		logger.Logger.Panicw("", err)
	}
	go router.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			logger.Logger.Sync()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
