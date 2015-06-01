//usr/local/go/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	"log"

	"github.com/jansplichal/commander/queue"
	"github.com/jansplichal/commander/ws"
)

const (
	amqpURL   = "amqp://guest:guest@localhost:5672/"
	queueName = "CMD_IN"
	wsPath    = "/"
	wsPort    = ":8888"
)

func main() {
	fmt.Println("Starting command Server")

	conn, msgs, err := queue.ListenQueue(amqpURL, queueName)

	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		panic("Can not connect to queue ...")
	}

	go func() {
		for d := range msgs {
			log.Printf("Reply to: %s", d.ReplyTo)
			log.Printf("Correlation id: %s", d.CorrelationId)
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	ws.Listen(wsPath, wsPort)
}
