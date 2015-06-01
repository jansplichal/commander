package ws

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

/*CmdRequest ...*/
type CmdRequest struct {
	recordID string
}

var clients = make(map[string]*websocket.Conn)
var requests = make(map[string]CmdRequest)

func webHandler(ws *websocket.Conn) {
	var err error

	//r, err := ioutil.ReadAll(ws.IsServerConn())
	client := ws.Request().Header.Get("X-Agent-Name")
	clients[client] = ws

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			delete(clients, client)
			fmt.Println("Can't receive")
			break
		}

		fmt.Println(ws == nil)

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

/*Listen ...*/
func Listen() {
	server := websocket.Server{
		Handler: websocket.Handler(webHandler),
	}

	http.Handle("/", server)
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
