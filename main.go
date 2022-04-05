package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
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

func SendMessage(mess []byte, queue string) {

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
			ContentType: "application/json",
			Body:        []byte(mess),
		})
	failOnError(err, "Failed to publish a message")
}

type User struct {
    Name string
		Age int32
}

type BasicTrade struct {
	Side         float64 	`json:"Side"`
	SecurityId   string 	`json:"SecurityId"`
	TransactTime string 	`json:"TransactTime"`
	AvgPrice     float64 	`json:"AvPrice"`
	ClOrdId      string 	`json:"ClOrdId"`
	OrderId      string 	`json:"OrderId"`
	Currency     string 	`json:"Currency"`
	Order_qty    float64	`json:"Order_qty"`
	ExecType     string 	`json:"ExecType"`
}

func main () {

	status := true
	for status {
		var mess string
		fmt.Println("Send order y/n (quit to exit): ")
		fmt.Scan(&mess)

		if mess == "quit" {
			status = false
			break
		} else {
			order := &BasicTrade{Side: 1,
				SecurityId:   "DE0007100000",
				TransactTime: "2022-03-29 14:47:27.698",
				ClOrdId:      "ec01b8fb-4301-4bd5-aab3-bedf45e42bb5",
				OrderId:      "63314325",
				Currency:     "EUR",
				Order_qty:    0.29,
				ExecType:     "A"}
			mess, _ := json.Marshal(order)
			SendMessage(mess, "orders")
		}
	}
}
