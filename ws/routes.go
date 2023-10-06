package ws

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ConnectWs(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		// vehicle_id := c.Query("vehicle_id")
		// vehicle_type := c.Query("vehicle_type")

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := NewClient(id, hub, conn)

		go client.writePump()
		go client.readPump()
	}
}
