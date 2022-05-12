package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"encoding/json"
	//"strconv"
	"math/rand"
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

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
			bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func SendMessage(mess []byte) {

	conn := ConnectToRabbitMQ("amqp://guest:guest@localhost:5672/")

	// Open a channel
	ch, err := conn.Channel()
	var failError string = "Failed to open a channel"
	failOnError(err, failError)
	defer ch.Close()

	corrId := randomString(32)

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	  )

	failOnError(err, "Failed to declare a queue")

	//msgs, err := ch.Consume(
	//		q.Name, // queue
	//		"",     // consumer
	//		true,   // auto-ack
	//		false,  // exclusive
	//		false,  // no-local
	//		false,  // no-wait
	//		nil,    // args
	//)
	failOnError(err, "Failed to register a consumer")
	  
	  err = ch.Publish(
		"",          // exchange
		"order-open", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
		  ContentType:   "text/plain",
		  CorrelationId: corrId,
		  ReplyTo:       q.Name,
		  Body:          []byte(mess),
	  })
}

type User struct {
    Name string
		Age int32
}

type IncomingTradeOrder struct {
	PartyID            string  `json:"party_id"`
	PartyRole          string  `json:"party_role"`
	PartySubIDType     string  `json:"party_sub_id_type"`
	OrderQty           float64 `json:"order_qty"`
	Side               string  `json:"side"`
	OrdType            string  `json:"ord_type"`
	OrderQuotationType string  `json:"order_quotation_type"`
	SecurityID         string  `json:"security_id"`
	ExDestination      string  `json:"ex_destination"`
	TimeInForce        int     `json:"time_in_force"`
	Currency           string  `json:"currency"`
	Quantity           float64 `json:"quantity"`
	ID                 string  `json:"id"`
	ExecType           string  `json:"exec_type"`
	Price              string  `json:"price"`
	TransactTime       string  `json:"transact_time"`
	PartySubID         string  `json:"party_sub_id"`
	PartyStammnummer   string  `json:"party_stammnummer"`
	CancelPending      string  `json:"cancel_pending"`
	Canceled           string  `json:"canceled"`
	ReplacePending     string  `json:"replace_pending"`
	Replaced           string  `json:"replaced"`
	OrderID            string  `json:"order_id"`
	ClOrdID            string  `json:"cl_ord_id"`
	CreatedAt          string  `json:"created_at"`
	OrdRejReason       string  `json:"ord_rej_reason"`
	RejectedText       string  `json:"rejected_text"`
	OrderReason        string  `json:"order_reason"`
	Cancellable        bool    `json:"cancellable"`
	PortfolioID        string  `json:"portfolio_id"`
	SubscriptionID     string  `json:"subscription_id"`
	TradingAccountID   string  `json:"trading_account_id"`
	Stamm              string  `json:"stamm"`
}

func main() {

	status := true
	for status {
		var mess string
		fmt.Println("Type quit to exit, enter to send order: ")
		fmt.Scan(&mess)

		if mess == "quit" {
			status = false
			break
		} else {
			order := &IncomingTradeOrder{PartyID: "BAADERBANK",
				PartyRole: "24",
				PartySubIDType: "22",
				OrderQty: 0,
				Side: "1",
				OrdType: "1",
				OrderQuotationType: "UNIT",
				SecurityID: "DE0007100000",
				ExDestination: "MUND",
				TimeInForce: 0,
				Currency: "EUR",
				Quantity: 0,
				ID: "d043819a-398a-4dca-afb4-6caa20fc22f8",
				ExecType: "",
				Price: "",
				TransactTime: "2022-05-12T10:44:38.095Z",
				PartySubID: "",
				PartyStammnummer: "",
				CancelPending: "",
				Canceled: "",
				ReplacePending: "",
				Replaced: "",
				OrderID: "",
				ClOrdID: "",
				CreatedAt: "2022-05-12T10:44:38.064Z",
				OrdRejReason: "",
				RejectedText: "",
				OrderReason: "",
				Cancellable: true,
				PortfolioID: "",
				SubscriptionID: "",
				TradingAccountID:"b53ae3c8-e0a9-40d2-bf95-876f56c48714",
				Stamm: "1000001"}
			mess, _ := json.Marshal(order)
			SendMessage(mess)
		}
	}
}
