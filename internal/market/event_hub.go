package market

import (
	"encoding/json"
	"sync"
)

// Event represents a real-time market event broadcast to SSE clients
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// EventHub manages Server-Sent Events (SSE) connections
type EventHub struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan Event
	mu         sync.RWMutex
}

func NewEventHub() *EventHub {
	hub := &EventHub{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan Event, 256),
	}
	go hub.run()
	return hub
}

func (h *EventHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
			}
			h.mu.Unlock()
		case event := <-h.broadcast:
			data, err := json.Marshal(event)
			if err != nil {
				continue
			}
			msg := string(data)
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client <- msg:
				default:
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *EventHub) Subscribe() chan string {
	ch := make(chan string, 64)
	h.register <- ch
	return ch
}

func (h *EventHub) Unsubscribe(ch chan string) {
	h.unregister <- ch
}

func (h *EventHub) Broadcast(event Event) {
	select {
	case h.broadcast <- event:
	default:
	}
}

func (h *EventHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
