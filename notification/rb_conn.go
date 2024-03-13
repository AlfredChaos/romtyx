package notification

import (
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	RabbitMQ   = "rabbitmq"
	RbAddress  = "0.0.0.0:5672"
	RbUser     = "guest"
	RbPassword = "guest"
	RbTopic    = "romtyx"
	RbLock     = sync.Mutex{}
	RbConn     *amqp.Connection
	RbChan     *amqp.Channel
	RbErr      error
)

func Connect(logger *logrus.Entry) error {
	RbLock.Lock()
	defer RbLock.Unlock()

	url := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		RbUser,
		RbPassword,
		RbAddress,
	)
	RbConn, RbErr = amqp.Dial(url)
	failOnError(RbErr, "Failed to connect to RabbitMQ", logger)

	return nil
}

func Disconnect() error {
	if RbConn != nil {
		if err := RbConn.Close(); err == nil {
			RbConn = nil
		} else {
			return err
		}
	}
	if RbChan != nil {
		if err := RbChan.Close(); err == nil {
			RbChan = nil
		} else {
			return err
		}
	}
	return nil
}

func failOnError(err error, msg string, logger *logrus.Entry) {
	if err != nil {
		logger.Panicf("%s: %s", msg, err)
	}
}
