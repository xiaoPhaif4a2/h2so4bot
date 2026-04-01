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
    // 构造 OneBot 11 标准的发送指令
    params := map[string]interface{}{
        "message": text,
    }

    // 根据消息类型决定是填 user_id 还是 group_id
    if event.MessageType == "private" {
        params["user_id"] = event.UserID
        params["message_type"] = "private"
    } else {
        params["group_id"] = event.GroupID
        params["message_type"] = "group"
    }

    msg := map[string]interface{}{
        "action": "send_msg", // 或者用 send_group_msg
        "params": params,
    }

    data, _ := json.Marshal(msg)
    conn.WriteMessage(websocket.TextMessage, data)
}