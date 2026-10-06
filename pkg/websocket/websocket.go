package websocket

import (
	"fmt"
	"net/http"

	wsocket "github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	token  string
	socket string
	conn   *wsocket.Conn
}

func NewWebsocketConnection(token, socket string) (*WebSocketConnection, error) {
	ws := &WebSocketConnection{
		token:  token,
		socket: socket,
	}

	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	conn, _, err := wsocket.DefaultDialer.Dial(socket, headers)
	if err != nil {
		return nil, err
	}

	authMessage := SocketMessage{
		Event: SocketEventTypeAuth,
		Args:  []string{token},
	}

	if err := conn.WriteJSON(authMessage); err != nil {
		return nil, err
	}

	for {
		var message SocketMessage
		err := conn.ReadJSON(&message)
		if err != nil {
			break
		}

		if message.Event == SocketEventTypeAuthSucess {
			break
		}
	}

	ws.conn = conn
	return ws, nil
}

func (wsc *WebSocketConnection) SendCommand(message SocketMessage) error {
	if wsc.conn == nil {
		return fmt.Errorf("ws connection is not ready yet")
	}

	if err := wsc.conn.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (wsc *WebSocketConnection) AuthAndListen(onEvent func(m SocketMessage)) error {
	for {
		var message SocketMessage
		err := wsc.conn.ReadJSON(&message)
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}

		onEvent(message)
	}

	return nil
}

func (wsc *WebSocketConnection) Shutdown() {
	if wsc.conn != nil {
		wsc.conn.Close()
	}
}
