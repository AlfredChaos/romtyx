package main

import (
	"flag"
	"os"
	"romtyx/database"

	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	logger := log.WithFields(log.Fields{"project": "migration"})

	commandPtr := flag.String("cmd", "", "up or down")
	flag.Parse()
	command := *commandPtr
	if command != "up" && command != "down" {
		logger.Errorf("Params not support!")
		return
	}

	migrationPath, _ := os.Getwd()
	if err := database.Connect(logger); err != nil {
		logger.Info("connect to database failed.")
	}
	sqlDb := database.DbConn.DB()

	goose.SetVerbose(true)
	if err := goose.SetDialect(database.MySQL); err != nil {
		logger.Errorf("set goose dialect %s error", database.MySQL)
		return
	}
	arguments := make([]string, 0)
	if err := goose.Run(command, sqlDb, migrationPath, arguments...); err != nil {
		logger.Errorf("migration occurs error: %v", err)
		return
	}
}
