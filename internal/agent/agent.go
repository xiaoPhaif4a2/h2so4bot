package agent

import (
	"fmt"
	"math/rand"
	"strings"
	"sync" // 🟢 引入并发锁

	"h2so4bot/config"
	"h2so4bot/internal/database"
	"h2so4bot/internal/llm"
	"h2so4bot/internal/onebot"
	"h2so4bot/internal/persona"

	"github.com/gorilla/websocket"
)

type Agent struct {
	Brain      *llm.Brain
	Config     *config.Config
	GroupMoods map[int64]*persona.MoodState // 🟢 每个群独立的情绪记录
	mu         sync.Mutex                   // 🟢 并发锁，防止多个群同时说话导致 map 崩溃
}

func NewAgent(brain *llm.Brain, cfg *config.Config) *Agent {
	return &Agent{
		Brain:      brain,
		Config:     cfg,
		GroupMoods: make(map[int64]*persona.MoodState), // 🟢 初始化 map
	}
}

// 🟢 新增一个安全获取或创建群情绪的方法
func (a *Agent) getGroupMood(groupID int64) *persona.MoodState {
	a.mu.Lock()         // 加锁
	defer a.mu.Unlock() // 函数结束时解锁

	mood, exists := a.GroupMoods[groupID]
	if !exists {
		// 如果这个群是第一次说话，给它新建一个情绪档案
		mood = persona.NewMoodState()
		a.GroupMoods[groupID] = mood
	}
	return mood
}

func (a *Agent) HandleMessage(conn *websocket.Conn, event onebot.MessageEvent) {
	database.SaveLog(event.UserID, event.GroupID, "user", event.RawMessage)

	shouldReply := false
	cleanedMsg := event.RawMessage
	atCode := fmt.Sprintf("[CQ:at,qq=%s]", a.Config.BotQQ)

	isPrivate := event.MessageType == "private"
	isAtMe := strings.Contains(event.RawMessage, atCode)

	// 🟢 获取当前群的专属情绪，并更新
	currentMood := a.getGroupMood(event.GroupID)
	currentMood.Update() // 不再传 isAtMe 参数了

	if isPrivate || isAtMe {
		shouldReply = true
		cleanedMsg = strings.ReplaceAll(event.RawMessage, atCode, "")
		cleanedMsg = strings.TrimSpace(cleanedMsg)
		if cleanedMsg == "" {
			cleanedMsg = "（只是默默地看着你）"
		}
	} else if event.MessageType == "group" {
		probability := 0.02

		if strings.Contains(cleanedMsg, "书") || strings.Contains(cleanedMsg, "文学") || strings.Contains(cleanedMsg, "名著") {
			probability += 0.15
		}

		// 🟢 根据该群的专属情绪决定接话欲
		switch currentMood.GetLevel() {
		case persona.MoodHappy:
			probability += 0.20 // 寂寞了，疯狂插话
		case persona.MoodAnnoyed:
			probability = 0.00
		}

		if rand.Float64() < probability {
			shouldReply = true
		}
	}

	if shouldReply {
		// 🟢 检查当前发消息的人（event.UserID）是不是配置里的主人（a.Config.MasterQQ）
		isCreator := event.UserID == a.Config.MasterQQ

		// 🟢 将判断结果传给人设生成器
		currentPersona := persona.BuildSystemPrompt(currentMood, isCreator)
		
		// 获取回复（此时短期记忆也会生效了）
		reply := a.Brain.Think(currentPersona, event.UserID, event.GroupID, cleanedMsg)

		if event.MessageID != 0 {
			reply = fmt.Sprintf("[CQ:reply,id=%d]%s", event.MessageID, reply)
		}

		database.SaveLog(event.UserID, event.GroupID, "assistant", reply)
		onebot.SendMessage(conn, event, reply)

		fmt.Printf(">> [群:%d|用户:%d] 情绪值: %d -> 🤖: %s\n", event.GroupID, event.UserID, currentMood.Value, reply)
	}
}
