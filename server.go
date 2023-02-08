package main

import "fmt"

type WsServer struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer(id string) *WsServer {
	return &WsServer{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)

		case message := <-server.broadcast:
			fmt.Println("Run")
			server.broadcastToClients(message)
		}
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	fmt.Println(len(server.clients))
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}
