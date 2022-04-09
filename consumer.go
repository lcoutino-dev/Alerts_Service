package main

import (
	"fmt"
	"os"

	"github.com/badkaktus/gorocket"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

/*
Load environment variables
*/
var rabbit_host = os.Getenv("RABBIT_HOST")
var rabbit_port = os.Getenv("RABBIT_PORT")
var rabbit_user = os.Getenv("RABBIT_USERNAME")
var rabbit_password = os.Getenv("RABBIT_PASSWORD")
var rabbit_queue_name = os.Getenv("RABBIT_QUEUE")
var rabbit_exchange_name = os.Getenv("RABBIT_EXCHANGE")
var rabbit_routing_key = os.Getenv("RABBIT_ROUTING_KEY")

func main() {
	consume()
}

func notifiRocketChat(msj string, client *gorocket.Client) bool {
	message := gorocket.Message{
		Channel: "PLD",
		Text:    msj,
	}
	msg, err := client.PostMessage(&message)
	if err != nil {
		fmt.Printf("Error: %+v", err)
	}
	return msg.Success
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func consume() {
	client := gorocket.NewClient(os.Getenv("ROCKET_SERVER"))
	// login as the main admin user
	login := gorocket.LoginPayload{
		User:     os.Getenv("ROCKET_USER"),
		Password: os.Getenv("ROCKET_PASSWORD"),
	}

	lg, err := client.Login(&login)

	if err != nil {
		fmt.Printf("Error: %+v", err)
	}
	fmt.Printf("Login to rocket chat is success. I'm %s \n", lg.Data.Me.Username)

	conn, err := amqp.Dial("amqp://" + rabbit_user + ":" + rabbit_password + "@" + rabbit_host + ":" + rabbit_port + "/")

	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
	}
	err = ch.ExchangeDeclare(
		rabbit_exchange_name, // name
		"direct",             // type
		true,                 // durable
		false,                // auto-deleted
		false,                // internal
		false,                // no-wait
		nil,                  // arguments
	)

	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		rabbit_queue_name, // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)

	if err != nil {
		log.Fatalf("%s: %s", "Failed to declare a queue", err)
	}

	fmt.Println("Channel and Queue established ✔️")

	err = ch.QueueBind(
		q.Name,               // queue name
		rabbit_routing_key,   // routing key
		rabbit_exchange_name, // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

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
			// log.Printf("Received a message: %s", d.Body)
			// fmt.Println(string(d.Body))
			// msjerror := string(d.Body)
			succes_message := notifiRocketChat(string(d.Body), client)
			if succes_message {
				d.Ack(false)
			}
		}
	}()

	fmt.Println("Recibiendo logs de \033[1m \033[93m KYCEVM... \033[0m 💀")
	<-forever
}
