package ws

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
	"os"
)

type WebSocketClient struct {
	ConfigStr string
	Sender    chan *Socket
}

var wsClient = &WebSocketClient{
	ConfigStr: fmt.Sprintf("%s:%s", os.Getenv("WEBSOCKET_HOST"), os.Getenv("WEBSOCKET_PORT")),
}

var sender = make(chan *Socket)
var ws *websocket.Conn

func GetWebsocketConnection() *WebSocketClient {
	return wsClient
}

func (conn *WebSocketClient) Send(channel, event string, data interface{}) {
	go ListenSenderChannel()

	sender <- &Socket{
		Channel: channel,
		Event:   event,
		Data:    data,
	}
}

func ListenSenderChannel() {
	for {
		select {
		case message, ok := <-sender:
			if !ok {
				continue
			}

			data, err := json.Marshal(message)
			if err != nil {
				slog.Error("Error marshalling message: ", err)
				continue
			}

			req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/ws/event", wsClient.ConfigStr), bytes.NewReader(data))
			if err != nil {
				slog.Error("Error creating request: ", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				slog.Error("Error sending request: ", err)
				continue
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				slog.Error("Error sending request: ", resp.Status)
				continue
			}
		}
	}
}
