package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second

	// EndOfTransmission end
	EndOfTransmission = "\u0004"
)

type WebTerminal struct {
	wsConn    *websocket.Conn
	sizeChan  chan remotecommand.TerminalSize
	doneChan  chan struct{}
	cancelCtx context.CancelFunc
	tty       bool
}

type TerminalMessage struct {
	Operation string `json:"operation"`
	Data      string `json:"data"`
	Rows      uint16 `json:"rows"`
	Cols      uint16 `json:"cols"`
}

func (t *WebTerminal) Next() *remotecommand.TerminalSize {
	select {
	case size := <-t.sizeChan:
		return &size
	case <-t.doneChan:
		return nil
	}
}

func (t *WebTerminal) Read(p []byte) (n int, err error) {
	_, message, err := t.wsConn.ReadMessage()
	if err != nil {
		return copy(p, EndOfTransmission), err
	}

	var msg TerminalMessage
	if err := json.Unmarshal([]byte(message), &msg); err != nil {
		// binary data receive
		return copy(p, message), nil
	}
	switch msg.Operation {
	case "stdin":
		return copy(p, []byte(msg.Data)), nil
	case "resize":
		t.sizeChan <- remotecommand.TerminalSize{Width: msg.Cols, Height: msg.Rows}
		return 0, nil
	case "ping":
		return 0, nil
	default:
		return copy(p, EndOfTransmission), fmt.Errorf("unknown message type '%s'", msg.Operation)
	}
}

func (t *WebTerminal) Write(p []byte) (n int, err error) {
	if err := t.wsConn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return copy(p, []byte(EndOfTransmission)), err
	}
	return len(p), nil
}

func (t *WebTerminal) Stdin() io.Reader {
	return t
}

func (t *WebTerminal) Stdout() io.Writer {
	return t
}

func (t *WebTerminal) Stderr() io.Writer {
	return t
}

func (t *WebTerminal) Done() {
	close(t.doneChan)
}

func (t *WebTerminal) Close() error {
	return t.wsConn.Close()
}

func (t *WebTerminal) Tty() bool {
	return t.tty
}
