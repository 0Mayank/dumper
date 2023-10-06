package ws

type LocationRoutine struct {
	from *Client
	to   *Client
	stop chan bool
	hub  *Hub
}

func newLocationRoutine(from, to *Client, stop chan bool, hub *Hub) *LocationRoutine {
	return &LocationRoutine{from, to, stop, hub}
}

func (l *LocationRoutine) run() {
	for {
		select {
		case <-l.stop:
			l.hub.unregisterRoutine <- l
			return
		case <-l.from.locationUpdate:
			l.from.locationLock.Lock()
			location := l.from.location
			l.from.locationLock.Unlock()

			msg := map[string]string{
				"type":      "location",
				"latitude":  location.latitude,
				"longitude": location.longitude,
			}

			l.to.send <- msg
		}
	}
}
