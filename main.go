package main

import (
	"fmt"
	"net/http"
	"os"
	"romtyx/api"
	"romtyx/database"
	"romtyx/notification"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var addr = "0.0.0.0:37011"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("Main Process Exception!")
			os.Exit(1)
		}
	}()

	logger := log.WithFields(log.Fields{"project": "romtyx"})

	start := time.Now()
	router := gin.New()
	registerRoutes(logger, router)
	registerDatabase(logger)
	registerNotification(logger)
	ser := &http.Server{Addr: addr, Handler: router}

	logger.Infof("server: listening on %s [%s]", addr, time.Since(start))
	if err := ser.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			// close database
			if dbErr := database.Disconnect(); dbErr != nil {
				logger.Error(dbErr)
			}
			logger.Infof("database closed.")
			// close rabbitmq
			if rbErr := notification.Disconnect(); rbErr != nil {
				logger.Error(rbErr)
			}
			logger.Infof("rabbitmq closed.")
			logger.Infof("server: shutdown complete")
		} else {
			logger.Errorf("server: %s", err)
		}
	}
}

func registerRoutes(logger *log.Entry, router *gin.Engine) {
	v1 := router.Group(fmt.Sprintf("%s/v1", addr))

	{
		api.Upload(logger, v1)
		api.Download(logger, v1)
	}
}

func registerDatabase(logger *log.Entry) {
	if err := database.Connect(logger); err != nil {
		logger.Error(err)
	}
}

func registerNotification(logger *log.Entry) {
	if err := notification.Connect(logger); err != nil {
		logger.Error(err)
	}
}
