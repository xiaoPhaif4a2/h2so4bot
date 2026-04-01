package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	APIKey   string
	WSURL    string
	MySQLDSN string
	BotQQ    string 
}

func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		APIKey:   os.Getenv("SILICON_API_KEY"),
		WSURL:    os.Getenv("QQ_BOT_WS_URL"),
		MySQLDSN: os.Getenv("MYSQL_DSN"),
		BotQQ:    os.Getenv("BOT_QQ"), 
	}
}