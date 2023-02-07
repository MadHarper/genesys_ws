package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var wsServer = NewWebsocketServer()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// Todo: переделать
	go wsServer.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("It's working!!!!!"))
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})

	address := fmt.Sprintf("0.0.0.0:%s", port)

	server := &http.Server{
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id, len(id))

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	// defer conn.Close()

	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer)

	go client.writePump()
	go client.readPump()

	// Todo передавать идентификатор чата?
	client.send <- []byte("first income")

	wsServer.register <- client
}
