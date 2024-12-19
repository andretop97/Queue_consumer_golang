package main

import (
	"fmt"
	"log"

	"github.com/andretop97/Queue_consumer_golang/src/consumer"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}


func batata(Body string) error {
	fmt.Println(Body)
	return nil
}

func main() {
	conn, err := amqp.Dial("amqp://exemple:1234@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	consumer, err := consumer.NewRabbitMQConsumer("amqp://exemple:1234@localhost:5672/", "teste")
	if err != nil{
		fmt.Print(err)
	}

	consumer.Consume(batata)
}