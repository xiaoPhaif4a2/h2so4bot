package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	APIKey string
	WSURL  string
}

func LoadConfig() *Config {
	godotenv.Load() // 加载 .env 文件
	return &Config{
		APIKey: os.Getenv("SILICON_API_KEY"),
		WSURL:  os.Getenv("QQ_BOT_WS_URL"),
	}
}