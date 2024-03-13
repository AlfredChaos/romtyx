package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Upload(logger *logrus.Entry, router *gin.RouterGroup) {
	router.POST("/upload", func(c *gin.Context) {

	})
}
