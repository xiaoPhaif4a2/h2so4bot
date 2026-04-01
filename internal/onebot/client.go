package onebot

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type MessageEvent struct {
	RawMessage  string `json:"raw_message"`
	MessageType string `json:"message_type"`
	UserID      int64  `json:"user_id"`
	GroupID     int64  `json:"group_id"`
	PostType    string `json:"post_type"`
}

// 封装发送逻辑，让外面调用更简单
func SendMessage(conn *websocket.Conn, event MessageEvent, text string) {
	msg := map[string]interface{}{
		"action": "send_msg",
		"params": map[string]interface{}{
			"message_type": event.MessageType,
			"user_id":      event.UserID,
			"group_id":     event.GroupID,
			"message":      text,
		},
	}
	data, _ := json.Marshal(msg)
	conn.WriteMessage(websocket.TextMessage, data)
}