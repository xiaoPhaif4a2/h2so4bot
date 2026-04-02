package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey       string
	WSURL        string
	MySQLDSN     string
	BotQQ        string
	TestGroupIDs []int64 // 存放多个测试群号
	MasterQQ     int64   // 🟢 确保这里是 int64 类型
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ 未找到 .env 文件，将尝试直接从系统环境变量读取")
	}

	// 1. 读取并解析多个群号
	groupIDsStr := os.Getenv("TEST_GROUP_IDS")
	var groupIDs []int64
	if groupIDsStr != "" {
		parts := strings.Split(groupIDsStr, ",")
		for _, part := range parts {
			cleanPart := strings.TrimSpace(part)
			if cleanPart == "" {
				continue
			}
			if parsedID, err := strconv.ParseInt(cleanPart, 10, 64); err == nil {
				groupIDs = append(groupIDs, parsedID)
			} else {
				log.Printf("⚠️ 群号解析失败，跳过无效值 [%s]: %v", cleanPart, err)
			}
		}
	}

	// 2. 🟢 读取并解析主人的 QQ 号
	masterQQStr := os.Getenv("MASTER_QQ")
	var masterQQ int64
	if masterQQStr != "" {
		parsedQQ, err := strconv.ParseInt(masterQQStr, 10, 64)
		if err == nil {
			masterQQ = parsedQQ
		} else {
			log.Printf("⚠️ MASTER_QQ 解析失败，请检查 .env 填写是否正确: %v", err)
		}
	}

	return &Config{
		APIKey:       os.Getenv("API_KEY"),
		WSURL:        os.Getenv("QQ_BOT_WS_URL"),
		MySQLDSN:     os.Getenv("MYSQL_DSN"),
		BotQQ:        os.Getenv("BOT_QQ"),
		TestGroupIDs: groupIDs,
		MasterQQ:     masterQQ, // 🟢 这里的变量必须和结构体里的类型（int64）一致
	}
}
