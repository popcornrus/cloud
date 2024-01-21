package main

import (
	"cloud/external/ws"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"log/slog"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all connections
}

var c *cache.Cache

func init() {
	c = cache.New(time.Hour*24, cache.NoExpiration)
}

var (
	conn *websocket.Conn
	err  error
)

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			slog.Error("Error closing connection: ", err)
		}
	}(conn)

	for {
		messageType, p, err := conn.ReadMessage()

		if err != nil {
			slog.Error("Error reading message:", err)
			break // Exit the loop if there's a read error
		}

		if messageType == websocket.TextMessage {
			var s *ws.Socket

			err = json.Unmarshal(p, &s)
			if err != nil {
				slog.Error("Error unmarshalling socket: ", err)
				continue
			}

			if s.Channel == "" {
				slog.Error("Channel is required")
				continue
			}

			if s.Event == "" {
				slog.Error("Event is required")
				continue
			}

			if s.Event[0:6] == "action" {
				ws.Actions(c, conn, s)
				continue
			}
		}
	}
}

func sendEvent(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Channel string      `json:"channel"`
		Event   string      `json:"event"`
		Data    interface{} `json:"data"`
	}

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var p request
	err := dec.Decode(&p)
	if err != nil {
		slog.Error("Error decoding request: ", err)
		return
	}

	data := &ws.Socket{
		Channel: p.Channel,
		Event:   p.Event,
		Data:    p.Data,
	}

	subscribedUsers, found := c.Get(data.Channel)
	if !found {
		slog.Error("Channel not found: ", data.Channel)
		return
	}

	for _, user := range subscribedUsers.([]*websocket.Conn) {
		err := user.WriteJSON(data)
		if err != nil {
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws/echo", handleConnection)
	http.HandleFunc("/ws/event", sendEvent)

	err := http.ListenAndServe(":5051", nil)
	if err != nil {
		return
	}
}
