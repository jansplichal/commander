//usr/local/go/bin/go run $0 $@ ; exit
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/jansplichal/commander/queue"
	"github.com/jansplichal/commander/ws"
)

const (
	amqpURL  = "amqp://guest:guest@localhost:5672/"
	queueIn  = "CMD_IN"
	queueOut = "CMD_OUT_ASYNC"
	wsPath   = "/"
	wsPort   = ":8888"
)

func main() {

	// con, err := db.Connect("localhost")
	// if err != nil {
	// 	panic("can not connect to database")
	// }
	// defer con.Close()
	//
	// con.Use("test", "users")
	//
	// err = con.Insert(&db.Person{Name: "Jan", Phone: "Splichal"}, &db.Person{Name: "Marie", Phone: "Kotrla"})
	// if err != nil {
	// 	fmt.Println("Can not insert person")
	// }
	//
	// r, err := con.FindOne(bson.M{"name": "Marie"}, &db.Person{})
	// // fmt.Println("Raw", p.Name)
	// if str, ok := r.(*db.Person); ok {
	// 	fmt.Println("Result from db", str.Name)
	// } else {
	// 	fmt.Println("Wrong conversion")
	// }

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	fmt.Println("Starting command Server")

	q, err := queue.ListenQueue(amqpURL, queueIn)

	if q.Conn != nil {
		defer q.Conn.Close()
	}

	if q.Channel != nil {
		defer q.Channel.Close()
	}

	if err != nil {
		panic("Can not connect to queue ...")
	}

	queue.PublishQueue(q)

	log.Printf(" [*] Waiting for messages on queue " + queueIn)
	go func() {
		ws.Listen(wsPath, wsPort)
	}()
	log.Printf(" [*] Web socket server initialized ...")

	<-done
}
