package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dsn string) {
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatalf("❌ 数据库 Ping 失败: %v", err)
	}
	log.Println("✅ 数据库连接成功")
}

func SaveLog(userID, groupID int64, role, content string) {
	_, err := DB.Exec("INSERT INTO chat_logs (user_id, group_id, role, content) VALUES (?, ?, ?, ?)",
		userID, groupID, role, content)
	if err != nil {
		log.Printf("⚠️ 保存聊天记录失败: %v", err)
	}
}

// 🟢 新增：用于给大脑传递历史消息的结构体
type LogMessage struct {
	Role    string
	Content string
}

// 🟢 新增：提取最近 N 条对话上下文
func GetContext(userID, groupID int64, limit int) []LogMessage {
	// 查询该群最近的聊天记录，按时间倒序（最新的在前面）
	rows, err := DB.Query("SELECT role, content FROM chat_logs WHERE group_id = ? ORDER BY id DESC LIMIT ?", groupID, limit)
	if err != nil {
		log.Printf("⚠️ 获取上下文失败: %v", err)
		return nil
	}
	defer rows.Close()

	var temp []LogMessage
	for rows.Next() {
		var msg LogMessage
		if err := rows.Scan(&msg.Role, &msg.Content); err == nil {
			temp = append(temp, msg)
		}
	}

	// ⚠️ 重点：反转顺序！因为大模型需要按时间正序阅读（最旧的在前，最新的在后）
	messages := make([]LogMessage, len(temp))
	for i, j := 0, len(temp)-1; i < len(temp); i, j = i+1, j-1 {
		messages[i] = temp[j]
	}

	return messages
}