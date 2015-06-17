//usr/local/go/bin/go msgun $0 $@ ; exit
package main

import (
	"encoding/json"
	"flag"

	"github.com/streadway/amqp"
)

const (
	url      = "amqp://guest:guest@localhost:5672/"
	queueIn  = "CMD_IN"
	queueOut = "CMD_OUT_ASYNC"
)

/*Message ...*/
type Message struct {
	Agent    string `json:"agent"`
	Module   string `json:"module"`
	Args     string `json:"args"`
	Machines string `json:"machines"`
}

func main() {
	con, err := amqp.Dial(url)
	defer con.Close()

	if err != nil {
		panic("No connection")
	}

	ch, err := con.Channel()
	defer ch.Close()

	if err != nil {
		panic("Cannot create channel")
	}

	msg, err := json.Marshal(Message{Agent: "agent1", Module: "ping", Args: "", Machines: "wccs"})

	n := flag.Uint("n", 1, "How many messages to put")
	flag.Parse()

	for i := 0; i < int(*n); i++ {
		ch.Publish("", queueIn, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	}
}
