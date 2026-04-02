package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"h2so4bot/config"
	"h2so4bot/internal/agent"
	"h2so4bot/internal/database"
	"h2so4bot/internal/llm"
	"h2so4bot/internal/onebot"

	"github.com/gorilla/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// 1. 加载配置
	cfg := config.LoadConfig()
	if cfg.APIKey == "" || cfg.MySQLDSN == "" || cfg.BotQQ == "" {
		log.Fatal("❌ 配置缺失：请检查 .env 是否填写了 API_KEY, MYSQL_DSN 和 BOT_QQ")
	}

	// 2. 初始化底层服务 (数据库和 LLM)
	database.InitDB(cfg.MySQLDSN)
	brain := llm.NewBrain(cfg.APIKey)

	// 3. 实例化智能体
	tsundereAgent := agent.NewAgent(brain, cfg)

	// 4. 连接 NapCat
	conn, _, err := websocket.DefaultDialer.Dial(cfg.WSURL, nil)
	if err != nil {
		log.Fatalf("❌ 无法拨通 NapCat: %v", err)
	}
	defer conn.Close()

	log.Printf("🌙 硫酸少女 (QQ:%s) 已在静谧中苏醒...", cfg.BotQQ)

	// --- 👇 这里是新增的主动插话（定时器）模块 👇 ---
	go func() {
		// ⚠️ 为了方便你立刻看到效果，我这里改成每 1 分钟检查一次。
		// 等测试成功后，你可以把它改回 1 * time.Hour (每小时)
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			nowHour := time.Now().Hour()
			// 确保只在白天活跃 (8点到23点)
			if nowHour >= 8 && nowHour <= 23 {
				// ⚠️ 为了测试能立刻触发，现在的概率是 1.0 (100%)。
				// 测试完毕后，建议改成 0.1 或 0.2
				if rand.Float64() < 1.0 {
					// 🟢 检查是否配置了群号
					if len(cfg.TestGroupIDs) == 0 {
						log.Println("⚠️ 未配置 TEST_GROUP_IDS，跳过本次主动发言")
						continue
					}

					// 🟢 随机挑选一个群
					randomIndex := rand.Intn(len(cfg.TestGroupIDs))
					targetGroupID := cfg.TestGroupIDs[randomIndex]

					log.Printf("⏳ 触发主动插话逻辑，准备在群 %d 发言，正在思考...", targetGroupID)

					sysPrompt := "你现在正独自在书房看书，心情平静。突然想对窗外的凡人们感叹一句关于命运或人性的文学名著句子，请直接输出感叹内容，带有你傲娇的语气，字数不超过30字。"
					msg := tsundereAgent.Brain.Think(sysPrompt, 0, targetGroupID, "（合上书本，望向窗外）")

					event := onebot.MessageEvent{
						GroupID:     targetGroupID,
						MessageType: "group",
					}
					onebot.SendMessage(conn, event, msg)
					log.Printf("🤖 主动发言 (群%d) -> %s\n", targetGroupID, msg)
				}
			}
		}
	}()
	// --- 👆 主动插话模块结束 👆 ---

	// 5. 开启事件循环（这里会一直卡住监听消息）
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("⚠️ WebSocket 连接中断")
			break
		}

		var event onebot.MessageEvent
		if err := json.Unmarshal(raw, &event); err != nil {
			continue
		}

		// 把收到的消息丢给 Agent 去处理
		if event.PostType == "message" {
			tsundereAgent.HandleMessage(conn, event)
		}
	}
}
