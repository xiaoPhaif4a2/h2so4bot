package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dsn string) {
	// dsn 格式: "user:password@tcp(127.0.0.1:3306)/mumubot"
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
}

func SaveLog(userID int64, groupID int64, role string, content string) {
	DB.Exec("INSERT INTO chat_logs (user_id, group_id, role, content) VALUES (?, ?, ?, ?)",
		userID, groupID, role, content)
}