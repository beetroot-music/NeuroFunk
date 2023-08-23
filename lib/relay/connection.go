package relay

import "github.com/gorilla/websocket"

type Connection struct {
	Socket *websocket.Conn
	Relay  *Relay
}

func (conn *Connection) Disconnect() (err error) {
	panic("not implemented")
}
