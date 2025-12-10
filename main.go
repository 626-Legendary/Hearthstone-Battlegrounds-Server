// main.go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"bgs-server/database" // ✅ 你的 database 包
)

func main() {
	// 1. 加载 .env
	_ = godotenv.Load()

	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("请在 .env 中配置 BLIZZARD_CLIENT_ID 和 BLIZZARD_CLIENT_SECRET")
	}

	// 2. 初始化数据库（建库 + 建 heroes 表）
	database.InitDB()

	log.Println("✅ 程序启动完成，数据库已初始化")
}
