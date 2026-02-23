package market

import (
	"encoding/json"
	"sync"
)

// Event represents a real-time market event broadcast to SSE clients
type Event struct {
	Type string      `json:"type"` // "trade", "price_update", "market_index"
	Data interface{} `json:"data"`
}

// EventHub manages Server-Sent Events (SSE) connections for real-time
// market data broadcasting. When a trade happens or price changes,
// the event is broadcast to all connected clients.
type EventHub struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan Event
	mu         sync.RWMutex
}

// NewEventHub creates and starts a new EventHub
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
					// Client buffer full, skip to prevent blocking
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Subscribe registers a new SSE client and returns its message channel
func (h *EventHub) Subscribe() chan string {
	ch := make(chan string, 64)
	h.register <- ch
	return ch
}

// Unsubscribe removes an SSE client
func (h *EventHub) Unsubscribe(ch chan string) {
	h.unregister <- ch
}

// Broadcast sends an event to all connected SSE clients
func (h *EventHub) Broadcast(event Event) {
	select {
	case h.broadcast <- event:
	default:
		// Broadcast channel full, skip to prevent blocking
	}
}

// ClientCount returns the number of connected SSE clients
func (h *EventHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
