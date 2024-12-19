package main

import (
	"errors"

	"github.com/andretop97/Queue_consumer_golang/src/consumer"
	"github.com/andretop97/Queue_consumer_golang/src/logger"
	"github.com/andretop97/Queue_consumer_golang/src/utils"
)

func handler(Body string) error {
	return errors.New("Só um erro")
}
func main() {
	logger.SetLoggerSettings()
	consumer, err := consumer.NewRabbitMQConsumer("amqp://exemple:1234@localhost:5672/", "teste")
	defer consumer.StopConsumer()

	utils.FailOnError(err, "Failed to create consumer")
	consumer.Consume(handler)
}