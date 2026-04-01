package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"h2so4bot/config"
	"h2so4bot/internal/database"
	"h2so4bot/internal/llm"
	"h2so4bot/internal/onebot"

	"github.com/gorilla/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	cfg := config.LoadConfig()
	if cfg.APIKey == "" || cfg.MySQLDSN == "" || cfg.BotQQ == "" {
		log.Fatal("❌ 配置缺失：请检查 .env 是否填写了 API_KEY, MYSQL_DSN 和 BOT_QQ")
	}

	database.InitDB(cfg.MySQLDSN)
	brain := llm.NewBrain(cfg.APIKey)

	conn, _, err := websocket.DefaultDialer.Dial(cfg.WSURL, nil)
	if err != nil {
		log.Fatalf("❌ 无法拨通 NapCat: %v", err)
	}
	defer conn.Close()

	log.Printf("🌙 硫酸少女 (QQ:%s) 已在静谧中苏醒...", cfg.BotQQ)

	for {
		// --- 这里是之前报错的地方，已经修改为 ReadMessage ---
		_, raw, err := conn.ReadMessage() 
		if err != nil {
			log.Println("⚠️ WebSocket 连接中断")
			break
		}

		var event onebot.MessageEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			continue
		}

		if event.PostType == "message" {
			processMessage(conn, brain, event, cfg)
		}
	}
}

func processMessage(conn *websocket.Conn, brain *llm.Brain, event onebot.MessageEvent, cfg *config.Config) {
	database.SaveLog(event.UserID, event.GroupID, "user", event.RawMessage)

	shouldReply := false
	cleanedMsg := event.RawMessage
	atCode := fmt.Sprintf("[CQ:at,qq=%s]", cfg.BotQQ)

	isPrivate := event.MessageType == "private"
	isAtMe := strings.Contains(event.RawMessage, atCode)

	if isPrivate || isAtMe {
		shouldReply = true
		cleanedMsg = strings.ReplaceAll(event.RawMessage, atCode, "")
		cleanedMsg = strings.TrimSpace(cleanedMsg)
		if cleanedMsg == "" {
			cleanedMsg = "（只是默默地看着你）"
		}
	} else if event.MessageType == "group" {
		probability := 0.05
		if strings.Contains(event.RawMessage, "书") || strings.Contains(event.RawMessage, "文学") || strings.Contains(event.RawMessage, "名著") {
			probability = 0.20
		}

		if rand.Float64() < probability {
			shouldReply = true
		}
	}

	if shouldReply {
		reply := brain.Think(event.UserID, cleanedMsg)

		database.SaveLog(event.UserID, event.GroupID, "assistant", reply)
		onebot.SendMessage(conn, event, reply)
		
		fmt.Printf(">> [%d] %s -> 🤖: %s\n", event.UserID, cleanedMsg, reply)
	}
}