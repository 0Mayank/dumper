package ws

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	maxMessageSize = 512
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
)

type Client struct {
	id             string
	hub            *Hub
	conn           *websocket.Conn
	send           chan map[string]string
	location       Coords
	locationLock   *sync.Mutex
	locationUpdate chan bool
}

func NewClient(id string, hub *Hub, conn *websocket.Conn) *Client {
	client := &Client{
		id:             id,
		hub:            hub,
		conn:           conn,
		send:           make(chan map[string]string, 10),
		location:       Coords{latitude: "", longitude: ""},
		locationLock:   &sync.Mutex{},
		locationUpdate: make(chan bool),
	}

	hub.registerClient <- client

	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregisterClient <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msg := map[string]string{}
		err := c.conn.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("[ERROR] %v", err)
			}
			break
		}

		switch msg["type"] {
		case "location":
			latitude, ok1 := msg["latitude"]
			longitude, ok2 := msg["longitude"]
			if !ok1 || !ok2 {
				writeError(c.conn, "latitude or longitude missing")
				continue
			}

			coords := Coords{latitude, longitude}

			c.locationLock.Lock()
			c.location = coords
			c.locationLock.Unlock()

			select {
			case c.locationUpdate <- true:
			default:
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteJSON(msg)
			if err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
