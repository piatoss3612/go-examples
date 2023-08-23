package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connectRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func declareQueue(ch *amqp.Channel, queueName string) error {
	_, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}

func publishMessage(ch *amqp.Channel, queueName, message string) error {
	err := ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	return err
}

func main() {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := createChannel(conn)
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := "my_queue"

	err = declareQueue(ch, queueName)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// publish messages to the queue
	for i := 1; i <= 10; i++ {
		err = publishMessage(ch, queueName, fmt.Sprintf("message: %d", i))
		if err != nil {
			log.Fatalf("Failed to publish a message: %v", err)
		}
		fmt.Printf("Published message %d\n", i)

	}
	fmt.Println("done publishing messages")

}
