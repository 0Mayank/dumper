package ws

import "log"

type Hub struct {
	clients           map[string]*Client
	registerClient    chan *Client
	unregisterClient  chan *Client
	broadcast         chan map[string]string
	locationRoutines  map[Pair]*LocationRoutine
	registerRoutine   chan *LocationRoutine
	unregisterRoutine chan *LocationRoutine
}

func NewHub() *Hub {
	return &Hub{
		clients:           make(map[string]*Client),
		registerClient:    make(chan *Client),
		unregisterClient:  make(chan *Client),
		broadcast:         make(chan map[string]string),
		locationRoutines:  make(map[Pair]*LocationRoutine),
		registerRoutine:   make(chan *LocationRoutine),
		unregisterRoutine: make(chan *LocationRoutine),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerClient:
			h.clients[client.id] = client
			log.Printf("New client registered: %s\n", client.id)

		case client := <-h.unregisterClient:
			log.Printf("Client unregistered: %s\n", client.id)
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}

		case routine := <-h.registerRoutine:
			log.Printf("New location routine registered: %s -> %s", routine.from.id, routine.to.id)
			h.locationRoutines[Pair{to: routine.to.id, from: routine.from.id}] = routine

		case routine := <-h.unregisterRoutine:
			log.Printf("Location routine unregistered: %s -> %s", routine.from.id, routine.to.id)
			if _, ok := h.locationRoutines[Pair{to: routine.to.id, from: routine.from.id}]; ok {
				delete(h.locationRoutines, Pair{to: routine.to.id, from: routine.from.id})
				close(routine.stop)
			}

		case message := <-h.broadcast:
			for id, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, id)
				}
			}
		}
	}
}
