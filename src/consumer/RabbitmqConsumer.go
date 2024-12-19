package consumer

import (
	"fmt"
	"log"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel *amqp.Channel
	queueName string
}

func (r *RabbitMQConsumer) createConnection(url string)  error{
	conn, err := amqp.Dial(url)

	if err != nil{
		return err
	}

	r.connection = conn
	return nil
}

func (r *RabbitMQConsumer) createChannel() error{
	channel, err := r.connection.Channel()

	if err != nil {
		return err
	}

	r.channel = channel
	return nil
}

func (r *RabbitMQConsumer) cretaeExchange(exchangeName string) error {
	err := r.channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) createQueue(queueName string, args amqp.Table) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) createBind(exchangeName string, queueName string) error {
	err := r.channel.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQConsumer) createExchangeAndQueueWithBind(exchangeName string, queueName string, args amqp.Table) error{
	err := r.cretaeExchange(exchangeName)
	if err != nil {
		return err
	}

	err = r.createQueue(queueName, args)
	if err != nil {
		return err
	}

	err = r.createBind(exchangeName, queueName)
	if err != nil {
		return err
	}

	return nil

}

func NewRabbitMQConsumer(url string, serviceName string) (*RabbitMQConsumer, error) {
	consumer := new(RabbitMQConsumer)
	err := consumer.createConnection(url)

	if err != nil {
		return nil, err
	}

	// defer consumer.connection.Close()

	err = consumer.createChannel()

	if err != nil {
		return nil, err
	}

	// defer consumer.channel.Close()

	exchangeName := fmt.Sprintf("%s-exchange", serviceName)
	queueName := fmt.Sprintf("%s-queue", serviceName)
	consumer.queueName = queueName
	
	// dlxName := fmt.Sprintf("%s-dlx", serviceName)
	// dlqName := fmt.Sprintf("%s-dlq", serviceName)

	queueConfig := make(amqp.Table)

	err = consumer.createExchangeAndQueueWithBind(exchangeName, queueName, queueConfig)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func (r *RabbitMQConsumer) getConsumer(queueName string) <-chan amqp.Delivery {
	delivery, err := r.channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	return delivery
}

func (r *RabbitMQConsumer) Consume(handler func(string) error){

	
	delivery := r.getConsumer(r.queueName)

	var forever chan struct{}

	go func() {
		for message := range delivery {
			body := string(message.Body)
			slog.Info("Mensagem recebida", "msg", body)
			err := handler(body)

			if err != nil {
				slog.Error("Erro ao consumir mensagem", "error", err)
				message.Nack(false, true)
				// message.Reject(true)
			}else{
				message.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	

}