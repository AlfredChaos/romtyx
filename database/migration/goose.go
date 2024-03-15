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

	commandPtr := flag.String("cmd", "", "up, down or create")
	namePtr := flag.String("name", "", "migration scripts name")
	stypePtr := flag.String("type", "", "migration scripts type")
	flag.Parse()
	command := *commandPtr
	name := *namePtr
	stype := *stypePtr
	arguments := make([]string, 0)
	if command != "up" && command != "down" && command != "create" {
		logger.Errorf("Params not support!")
		return
	}
	if command == "create" {
		if name == "" {
			logger.Errorf("name required")
			return
		}
		arguments = append(arguments, name)
		if stype != "sql" && stype != "go" {
			logger.Errorf("type unsupport")
			return
		}
		arguments = append(arguments, stype)
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

	if err := goose.Run(command, sqlDb, migrationPath, arguments...); err != nil {
		logger.Errorf("migration occurs error: %v", err)
		return
	}
}
