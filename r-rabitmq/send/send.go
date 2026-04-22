package main

import (
	"context"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// func FailOnError(err error, msg string) {
// 	if err != nil {
// 		fmt.Println(err)
// 		log.Panicf("%s: %s", msg, err)
// 	}
// }

func SendMessage() {

	url := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp091.Dial(url)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	randomUserID := time.Now().Unix()
	body := "Hello World! " + string(randomUserID)

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})

	FailOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)

}

// func main() {
// 	SendMessage()
// }
