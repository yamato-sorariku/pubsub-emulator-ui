package controller

import (
	"cloud.google.com/go/pubsub"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan WsMessage)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsMessage struct {
	Message string `json:"message"`
}

// クライアントからは JSON 形式で受け取る
type PubSubBody struct {
	Message      PubSubMessage `json:"message"`
	Subscription string        `json:"subscription"`
}

type PubSubMessage struct {
	Data      []byte `json:"data"`
	MessageId string `json:"messageId"`
}

type WsPostMessage struct {
	Data      string `json:"data"`
	MessageId string `json:"messageId"`
}

func HandleClients(c *gin.Context) {
	w := c.Writer
	r := c.Request
	go broadcastMessagesToClients()
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	defer socket.Close()

	clients[socket] = true

	for {
		var message WsMessage
		err := socket.ReadJSON(&message)
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			delete(clients, socket)
			break
		}
		broadcast <- message
	}
}

func BroadcastMessagesToClients(meg *pubsub.Message) {
	var wsPostMessage WsPostMessage
	wsPostMessage.MessageId = meg.ID
	wsPostMessage.Data = string(meg.Data)

	jsonData, err := json.Marshal(wsPostMessage)
	if err != nil {
		log.Printf("PubSub Json parse Error: %v", err)
	}
	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("error occurred while writing message to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func PingToClients() {
	for client := range clients {
		err := client.WriteMessage(websocket.PingMessage, []byte("ping"))
		if err != nil {
			log.Printf("error occurred while writing message to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func broadcastMessagesToClients() {
	for {
		message := <-broadcast
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error occurred while writing message to client: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
