package database

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

var (
	MySQL      = "mysql"
	DbAddress  = "tcp(0.0.0.0:3306)"
	DbUser     = "root"
	DbPassword = "root"
	DbName     = "romtyx"
	DbLock     = sync.Mutex{}
	DbConn     *gorm.DB
	DbConnErr  error
)

func Connect(logger *logrus.Entry) error {
	DbLock.Lock()
	defer DbLock.Unlock()

	driver := MySQL
	logger.Infof("Get database driver: %s", driver)
	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
		DbUser,
		DbPassword,
		DbAddress,
		DbName,
	)
	logger.Infof("Get database dns: %s", dsn)

	if driver == "" {
		return errors.New("database driver not specified")
	}
	if dsn == "" {
		return errors.New("database dsn not specified")
	}

	DbConn, DbConnErr = gorm.Open(driver, dsn)
	if DbConnErr != nil || DbConn == nil {
		// retry
		for i := 1; i <= 12; i++ {
			DbConn, DbConnErr = gorm.Open(driver, dsn)
			if DbConn != nil && DbConnErr == nil {
				break
			}
			time.Sleep(5 * time.Second)
		}

		if DbConnErr != nil || DbConn == nil {
			return DbConnErr
		}
	}

	DbConn.LogMode(true)
	DbConn.SetLogger(logger)

	maxConns := databaseMaxConns(0)
	DbConn.DB().SetMaxOpenConns(maxConns)
	DbConn.DB().SetMaxIdleConns(databaseMaxIdle(0, maxConns))
	DbConn.DB().SetConnMaxLifetime(time.Hour)

	if err := checkDbConn(logger, DbConn); err != nil {
		logger.Error("connect database error")
		return err
	}
	logger.Info("connect database success.")

	return nil
}

func Disconnect() error {
	if DbConn != nil {
		if err := DbConn.Close(); err == nil {
			DbConn = nil
			DbConnErr = nil
		} else {
			return err
		}
	}
	return nil
}

func DbClient() *gorm.DB {
	return DbConn
}

func databaseMaxConns(num int) int {
	limit := 0
	if num <= 0 {
		limit = (runtime.NumCPU() * 2) + 16
	}
	if limit > 1024 {
		limit = 1024
	}
	return limit
}

func databaseMaxIdle(num, conns int) int {
	limit := 0
	if num <= 0 {
		limit = runtime.NumCPU() + 8
	}
	if limit > conns {
		limit = conns
	}
	return limit
}

func checkDbConn(logger *logrus.Entry, db *gorm.DB) error {
	type Res struct {
		Value string `gorm:"column:Value;"`
	}

	var res Res
	if err := db.Raw("SHOW VARIABLES LIKE 'innodb_version'").Scan(&res).Error; err != nil {
		return nil
	} else if v := strings.Split(res.Value, "."); len(v) < 3 {
		logger.Warnf("config: unknown database server version")
	}
	return nil
}

type ModelBase struct {
	ID        int64     `gorm:"type:BIGINT;primary_key;" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime" json:"updated_at"`
}

func (mb *ModelBase) BeforeCreate(tx *gorm.DB) {

}
