package main

import (
	"fmt"
	"time"
)

type Hub struct {
	servers map[string]*WsServer
}

func NewHub() *Hub {
	return &Hub{
		servers: make(map[string]*WsServer),
	}
}

func (h *Hub) start() {
	go h.removeEmptyServers()
}

func (h *Hub) addServer(s *WsServer) error {
	if val := h.servers[s.id]; val != nil {
		return fmt.Errorf("uuid for server already exists in the Hub")
	}

	h.servers[s.id] = s
	return nil
}

func (h *Hub) findServer(id string) (*WsServer, error) {
	val := h.servers[id]
	if val == nil {
		return nil, fmt.Errorf("server not found in the Hub")
	}

	return val, nil
}

func (h *Hub) removeEmptyServers() {
	ticker := time.NewTicker(120 * time.Second)

	for {
		select {
		case <-ticker.C:
			for idx, server := range h.servers {
				if len(server.clients) == 0 {
					delete(h.servers, idx)
				}
			}
		}
	}
}
