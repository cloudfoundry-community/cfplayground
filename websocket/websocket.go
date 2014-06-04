package websocket

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type Pipe struct {
	In     chan []byte
	Out    chan *Message
	Prompt chan []byte
}

type Message struct {
	Cmd     string
	MsgType string
	Msg     string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func New(w http.ResponseWriter, r *http.Request) (*Pipe, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return &Pipe{}, errors.New("Failed to upgrade websocket connection: " + err.Error())
	} else {
		pipe := &Pipe{
			make(chan []byte, 0),
			make(chan *Message, 0),
			make(chan []byte, 0),
		}
		pipe.startStream(conn)
		return pipe, nil
	}
}

func (p *Pipe) startStream(conn *websocket.Conn) {
	go func() {
		defer func() {
			conn.Close()
		}()
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			} else {
				p.In <- msg
			}
		}
	}()

	go func() {
		var outMsg *Message
		for {
			outMsg = <-p.Out
			b, _ := json.Marshal(outMsg)
			conn.WriteMessage(websocket.TextMessage, b)
		}
	}()

}
