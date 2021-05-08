package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	body := bodyFrom(os.Args)
	err = ch.Publish(
		"logs_topic",
		serviceFrom(os.Args)+"."+severityFrom(os.Args),
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[3] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[3:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "info"
	} else {
		s = os.Args[2]
	}
	return s
}

func serviceFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "financeiro"
	} else {
		s = os.Args[1]
	}
	return s
}
