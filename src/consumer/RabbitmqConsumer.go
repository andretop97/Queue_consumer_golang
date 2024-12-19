package consumer

import (
	"fmt"
	"log/slog"

	"github.com/andretop97/Queue_consumer_golang/src/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

type queueConfig struct {
	dlxName *string
	msgTTL *int64
}

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel *amqp.Channel
	queueName string
}

func (r *RabbitMQConsumer) createConnection(url string){
	conn, err := amqp.Dial(url)
	utils.FailOnError(err, "Failed to initialize connection")
	r.connection = conn
}

func (r *RabbitMQConsumer) createChannel(){
	channel, err := r.connection.Channel()
	utils.FailOnError(err, "Failed to create channel")
	r.channel = channel
}

func (r *RabbitMQConsumer) cretaeExchange(exchangeName string) {
	err := r.channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare Exchange")
}

func (r *RabbitMQConsumer) createQueue(queueName string, queueConfig queueConfig) {
	args := amqp.Table{}

	if queueConfig.dlxName != nil {
		args["x-dead-letter-exchange"] = *queueConfig.dlxName
	}

	if queueConfig.msgTTL != nil {
		args["x-message-ttl"] = *queueConfig.msgTTL
	}

	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		args,
	)

	utils.FailOnError(err, "Failed to declare queue")
}

func (r *RabbitMQConsumer) createBind(exchangeName string, queueName string) {
	err := r.channel.QueueBind(
		queueName,
		"",
		exchangeName,
		false,
		nil,
	)

	utils.FailOnError(err, "Failed to create bind between exchange and queue")
}

func (r *RabbitMQConsumer) createExchangeAndQueueWithBind(exchangeName string, queueName string, queueConfig queueConfig) {
	r.cretaeExchange(exchangeName)
	slog.Info(fmt.Sprintf("Exchange %s created successfuly", exchangeName))
	r.createQueue(queueName, queueConfig)
	slog.Info(fmt.Sprintf("Queue %s created successfuly", queueName))
	r.createBind(exchangeName, queueName)
	slog.Info("Bind created successfuly")

}

func (consumer *RabbitMQConsumer) createSimpleQueue(serviceName string){
	exchangeName := fmt.Sprintf("%s-exchange", serviceName)
	queueName := fmt.Sprintf("%s-queue", serviceName)
	consumer.queueName = queueName
	consumer.createExchangeAndQueueWithBind(exchangeName, queueName, queueConfig{})

}

func (consumer *RabbitMQConsumer) createQueueWithDlq(serviceName string){
	exchangeName := fmt.Sprintf("%s-exchange", serviceName)
	queueName := fmt.Sprintf("%s-queue", serviceName)

	dlxName := fmt.Sprintf("%s-dlx", serviceName)
	dlqName := fmt.Sprintf("%s-dlq", serviceName)

	qConfig := queueConfig{}
	qConfig.dlxName = &dlxName

	
	consumer.createExchangeAndQueueWithBind(dlxName, dlqName, queueConfig{})

	consumer.createExchangeAndQueueWithBind(exchangeName, queueName, qConfig)

}

func NewRabbitMQConsumer(url string, serviceName string) (*RabbitMQConsumer, error) {
	consumer := new(RabbitMQConsumer)
	consumer.createConnection(url)
	consumer.createChannel()

	consumer.createQueueWithDlq(serviceName)

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
	utils.FailOnError(err, "Failed to create consumer")
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
				message.Nack(false, false)
				// message.Reject(true)
			}else{
				message.Ack(false)
			}
		}
	}()

	slog.Info(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	

}

func (r *RabbitMQConsumer) StopConsumer(){
	err := r.channel.Close(); 
	utils.FailOnError(err, "Filed to close channel")
    err = r.connection.Close()
	utils.FailOnError(err, "Filed to close connection")
}