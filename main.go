package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

var hub = NewHub()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	hub.start()

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		index(w, r)
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
	var wsServer *WsServer
	var err error
	var newId string

	if id == "" {
		newId = uuid.New().String()
		wsServer = NewWebsocketServer(newId)

		err = hub.addServer(wsServer)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		go wsServer.Run()
	} else {
		wsServer, err = hub.findServer(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, wsServer)

	go client.writePump()
	go client.readPump()
	wsServer.register <- client

	if newId != "" {
		event := fmt.Sprintf("{\"event\":\"serverCreated\",\"id\":\"%s\"}", newId)
		client.send <- []byte(event)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	lineBreak := "\r\n"
	info := "Данные по ws серверам:" + lineBreak

	for idx, val := range hub.servers {
		info = info + "Сервер " + idx + " | Клиентов:  " + strconv.Itoa(len(val.clients)) + lineBreak
	}

	w.Write([]byte(info))
}
