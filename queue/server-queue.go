package queue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

/*Serve ...*/
func Serve() {
	conn, msgs, err := listenQueue("amqp://guest:guest@localhost:5672/", "CMD_IN")

	if err != nil {
		panic("Can not connect to queue ...")
	}

	if conn != nil {
		defer conn.Close()
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	var wait string
	fmt.Scanln(&wait)

}

func listenQueue(url, queueName string) (*amqp.Connection, <-chan amqp.Delivery, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return conn, nil, err
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return conn, nil, err
	}

	err = ch.Qos(
		3,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return conn, nil, err
	}

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
		return conn, nil, err
	}
	return conn, msgs, nil
}
