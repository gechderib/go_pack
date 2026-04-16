package main

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func ReceiveMessage() {

	// Create Connection
	url := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp091.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// create channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// declare queue, its idempotent, will only be created if it doesn't exist already
	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// consume messages
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// ack	the message after processing
			d.Ack(false) // Ack only this one message
			// d.Ack(true)  // Ack this and all previous unacknowledged messages

			// // nack the message, requeue it for redelivery
			// d.Nack(false, true) // Nack only this one message and requeue it
			// d.Nack(true, true)  // Nack this and all previous unacknowledged messages and requeue them

			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func main() {
	ReceiveMessage()
}
