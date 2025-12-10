package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"

	"bgs-server/blizzard"
	"bgs-server/database"
)

func main() {
	_ = godotenv.Load()

	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("请在 .env 中配置 BLIZZARD_CLIENT_ID 和 BLIZZARD_CLIENT_SECRET")
	}

	// 1. 初始化数据库（建库 + heroes 表）
	database.InitDB()

	// 2. 拿 token
	tokenResp, err := blizzard.GetAccessToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("获取 access_token 失败: %v", err)
	}

	// 3. 自动翻页获取所有英雄
	heroes, err := blizzard.GetBattlegroundHero(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取战棋英雄失败: %v", err)
	}
	log.Printf("准备写入英雄数量: %d\n", len(heroes))

	// 4. 写入数据库：hs_id 存在则更新，不存在则创建
	db := database.DB

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en",
				"name_zh",
				"armor",
				"hero_power_id",
				"companion_id",
				"image_en",
				"image_zh",
				"is_duo",
				"is_solo",
			}),
		},
	).Create(&heroes).Error
	if err != nil {
		log.Fatalf("写入<英雄数据>到数据库失败: %v", err)
	}

	log.Println("✅ 酒馆战旗<英雄数据>数据库同步完成")
}
