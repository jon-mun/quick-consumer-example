package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var rabbit_host string
var rabbit_port string
var rabbit_user string
var rabbit_password string

func main() {
	godotenv.Load()

	rabbit_host = os.Getenv("RABBIT_HOST")
	rabbit_port = os.Getenv("RABBIT_PORT")
	rabbit_user = os.Getenv("RABBIT_USERNAME")
	rabbit_password = os.Getenv("RABBIT_PASSWORD")

	consume()
}

func consume() {

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}

	q, err := ch.QueueDeclare(
		"publisher", // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Channel and Queue established")

	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			d.Ack(false)
		}
	}()

	fmt.Println("Running...")
	<-forever
}
