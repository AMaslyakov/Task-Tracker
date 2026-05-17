package events

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type SSEHub struct {
	clients    map[chan Event]bool
	broadcast  chan Event
	register   chan chan Event
	unregister chan chan Event
	mu         sync.RWMutex
}

func NewSSEHub() *SSEHub {
	return &SSEHub{
		clients:    make(map[chan Event]bool),
		broadcast:  make(chan Event, 100),
		register:   make(chan chan Event),
		unregister: make(chan chan Event),
	}
}

func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("SSE client connected, total: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client)
				log.Printf("SSE client disconnected, total: %d", len(h.clients))
			}
			h.mu.Unlock()

		case event := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client <- event:
				default:
					close(client)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *SSEHub) Broadcast(event Event) {
	select {
	case h.broadcast <- event:
	default:
		log.Printf("SSE broadcast channel full, dropping event: %s", event.EventType)
	}
}

func (h *SSEHub) SSEHandler(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	client := make(chan Event, 10)
	h.register <- client

	defer func() {
		h.unregister <- client
	}()

	notify := c.Request.Context().Done()

	for {
		select {
		case <-notify:
			log.Println("SSE: client context cancelled")
			return
		case c := <-c.Writer.CloseNotify():
			log.Printf("SSE: client closed connection: %v", c)
			return
		case event, ok := <-client:
			if !ok {
				return
			}
			data, err := json.Marshal(map[string]interface{}{
				"event_type": event.EventType,
				"payload":    event.Payload,
			})
			if err != nil {
				log.Printf("failed to marshal SSE data: %v", err)
				continue
			}
			_, err = c.Writer.Write([]byte(
				fmt.Sprintf("event: %s\ndata: %s\n\n", event.EventType, data),
			))
			if err != nil {
				log.Printf("failed to write SSE: %v", err)
				return
			}
			c.Writer.Flush()
		}
	}
}
