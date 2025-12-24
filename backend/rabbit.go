package backend

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Function used to create a new RabbitMQ connection and channel
func createRabbitChannel() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := channel.QueueDeclare(
		"analytics_queue", // name
		true,              // durable
		false,             // autoDelete
		false,             // exclusive
		false,             // noWait
		nil,               // args
	)
	failOnError(err, "Failed to declare a queue")

	return conn, channel, &q
}

// Function used to create a new RabbitMQ connection and channel
func PublishRabbitMessage(body string) error {
	_, channel, queue := createRabbitChannel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(
		ctx,
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")

	return nil
}

// Server that listens for messages from RabbitMQ
func RabbitConsumerServer() {
	conn, channel, queue := createRabbitChannel()

	defer conn.Close()
	defer channel.Close()

	msgs, err := channel.Consume(
		queue.Name,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
