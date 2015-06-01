//usr/local/go/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	"log"

	"github.com/jansplichal/commander/queue"
	"github.com/jansplichal/commander/ws"
)

func main() {
	fmt.Println("Starting command Server")

	conn, msgs, err := queue.ListenQueue("amqp://guest:guest@localhost:5672/", "CMD_IN")

	if conn != nil {
		defer conn.Close()
	}

	if err != nil {
		panic("Can not connect to queue ...")
	}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	ws.Listen()
}
