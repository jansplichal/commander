//usr/local/go/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2/bson"

	"github.com/jansplichal/commander/db"
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

	con, err := db.Connect("localhost")
	if err != nil {
		panic("can not connect to database")
	}
	defer con.Close()

	con.Use("test", "users")

	err = con.Insert(&db.Person{Name: "Jan", Phone: "Splichal"}, &db.Person{Name: "Marie", Phone: "Kotrla"})
	if err != nil {
		fmt.Println("Can not insert person")
	}

	r, err := con.FindOne(bson.M{"name": "Marie"}, &db.Person{})
	// fmt.Println("Raw", p.Name)
	if str, ok := r.(*db.Person); ok {
		fmt.Println("Result from db", str.Name)
	} else {
		fmt.Println("Wrong conversion")
	}

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
