package notification

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	RabbitMQ     = "rabbitmq"
	RbAddress    = "0.0.0.0:5672"
	RbUser       = "admin"
	RbPassword   = "admin"
	RbTopic      = "exchange.romtyx.direct"
	RbRoutingKey = "normal"
	RbLock       = sync.Mutex{}
	RbConn       *amqp.Connection
	RbChan       *amqp.Channel
	RbErr        error
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
	logger.Infof("connect to RabbitMQ success.")

	// open channel
	RbChan, RbErr = RbConn.Channel()
	failOnError(RbErr, "Failed to open a channel", logger)
	logger.Info("open channel success.")

	// declare exchange
	RbErr = RbChan.ExchangeDeclare(
		RbTopic,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(RbErr, "Failed to declare an exchange", logger)
	logger.Infof("declare %s exchange success.", RbTopic)

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

func Publish(logger *logrus.Entry, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	RbErr = RbChan.PublishWithContext(
		ctx,
		RbTopic,
		RbRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	failOnError(RbErr, "Failed to publish a message", logger)
	logger.Infof("message has Sent %s", message)
}

func failOnError(err error, msg string, logger *logrus.Entry) {
	if err != nil {
		logger.Panicf("%s: %s", msg, err)
	}
}
