package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Download(logger *logrus.Entry, router *gin.RouterGroup) {
	router.GET("/download", func(c *gin.Context) {

	})
}
