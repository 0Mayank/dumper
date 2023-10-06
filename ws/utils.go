package ws

import "github.com/gorilla/websocket"

func writeError(conn *websocket.Conn, e string) {
	errJson := map[string]string{
		"error": e,
	}
	conn.WriteJSON(errJson)
}
