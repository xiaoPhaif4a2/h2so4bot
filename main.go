package main

import (
	"encoding/json"
	"myownbot/config"
	"myownbot/internal/llm"
	"myownbot/internal/onebot"
	"github.com/gorilla/websocket"
)

func main() {
	cfg := config.LoadConfig()
	brain := llm.NewBrain(cfg.APIKey)

	// 连接 QQ
	conn, _, _ := websocket.DefaultDialer.Dial(cfg.WSURL, nil)

	for {
		_, raw, _ := conn.ReadMessage()
		var event onebot.MessageEvent
		json.Unmarshal(raw, &event)

		if event.PostType == "message" {
			reply := brain.Think(event.RawMessage)
			onebot.SendMessage(conn, event, reply)
		}
	}
}