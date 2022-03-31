package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConnectToRabbitMQ(path string) *amqp.Connection {
	conn, err := amqp.Dial(path)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	return conn
}

func SendMessage(mess string, queue string) {

	conn := ConnectToRabbitMQ("amqp://guest:guest@localhost:5672/")

	// Open a channel
	ch, err := conn.Channel()
	var failError string = "Failed to open a channel"
	failOnError(err, failError)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//body := "Hello World!"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mess),
		})
	failOnError(err, "Failed to publish a message")
}

func main() {
	var status bool = true
	for status {
		var mess string
		fmt.Println("Input your message (quit to exit): ")
		fmt.Scan(&mess)
		if mess == "quit" {
			status = false
			break
		} else {
			SendMessage(mess, "nyqueue")
		}
	}
}
