package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

type Worker struct {
	id     int
	taskCh chan string
	wg     *sync.WaitGroup
}

func (w *Worker) run() {
	log.Printf("Worker [%v] running", w.id)

	for task := range w.taskCh {
		log.Printf("Worker [%d] processing task: (%s)", w.id, task)
	}

	log.Printf("Worker [%d] stopped\n", w.id)

	w.wg.Done() // signal that this worker has finished
}

func main() {
	var (
		queueName = "my_queue"
	)
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Successfully connected to RabbitMQ instance")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	err = declareQueue(channel, queueName)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	var (
		TOTAL_WORKERS = 5
		taskCh        = make(chan string, 10)
		waitGroup     = &sync.WaitGroup{}
	)

	// start worker pool
	for i := 0; i < TOTAL_WORKERS; i++ {
		worker := &Worker{
			id:     i,
			taskCh: taskCh,
			wg:     waitGroup,
		}
		waitGroup.Add(1) // increment the WaitGroup for each worker
		go worker.run()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

Loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("All messages consumed")
			break Loop
		case msg := <-msgs:
			taskCh <- string(msg.Body)
		}
	}

	close(taskCh)

	// wait for all workers to finish before exiting
	waitGroup.Wait()

	fmt.Println("Successfully consumed message from RabbitMQ")
}
