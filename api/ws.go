package api

import (
	"github.com/gorilla/websocket"
)

type DataMessage struct {
	Data any `json:"data"`
}

type ErrorMessage struct {
	Error error `json:"error"`
}

func SendErrorMessage(conn *websocket.Conn, err error) {
	conn.WriteJSON(ErrorMessage{err})
}

func SendDataMessage(conn *websocket.Conn, data any) {
	conn.WriteJSON(DataMessage{data})
}
