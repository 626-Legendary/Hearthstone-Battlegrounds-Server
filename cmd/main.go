package main

import (
	"log"
	"os"

	"bgs-server/blizzard"
	"bgs-server/database"
	"bgs-server/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm/clause"
)

const (
	AppName   = "Hearthstone-Battlegrounds-Server"
	Author    = "ZEXIANG ZHANG"
	Version   = "v0.0.2"
	BuildDate = "2025-12-10"
)

func main() {
	// 读取 .env
	_ = godotenv.Load()

	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("请在 .env 中配置 BLIZZARD_CLIENT_ID 和 BLIZZARD_CLIENT_SECRET")
	}

	// 1. 初始化数据库
	database.InitDB()
	db := database.DB

	// 2. 获取 token
	tokenResp, err := blizzard.GetAccessToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("获取 access_token 失败: %v", err)
	}

	// ======================================================
	// 3. 初始化 Classes
	// ======================================================
	classes := models.GetClasses()
	log.Printf("准备写入职业数量: %d\n", len(classes))

	err = db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name_en", "name_zh"}),
		},
	).Create(&classes).Error
	if err != nil {
		log.Fatalf("写入 classes 失败: %v", err)
	}
	log.Println("✅ 职业/种族表初始化完成")

	// ======================================================
	// 4. 英雄
	// ======================================================
	heroes, err := blizzard.GetHeroCards(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取战棋英雄失败: %v", err)
	}
	log.Printf("准备写入英雄数量: %d\n", len(heroes))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh", "armor", "hero_power_id",
				"companion_id", "image_en", "image_zh",
				"is_duo", "is_solo",
			}),
		},
	).Create(&heroes).Error
	if err != nil {
		log.Fatalf("写入英雄失败: %v", err)
	}
	log.Println("✅ 英雄同步完成")

	// ======================================================
	// 5. 大饰品
	// ======================================================
	greaterTrinkets, err := blizzard.GetGreaterTrinketsCards(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取大饰品失败: %v", err)
	}
	log.Printf("准备写入大饰品数量: %d\n", len(greaterTrinkets))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh", "mana_cost",
				"text_en", "text_zh",
				"image_en", "image_zh",
				"trinket_type",
			}),
		},
	).Create(&greaterTrinkets).Error
	if err != nil {
		log.Fatalf("写入大饰品失败: %v", err)
	}

	// ======================================================
	// 6. 小饰品
	// ======================================================
	lesserTrinkets, err := blizzard.GetLesserTrinketsCards(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取小饰品失败: %v", err)
	}
	log.Printf("准备写入小饰品数量: %d\n", len(lesserTrinkets))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh", "mana_cost",
				"text_en", "text_zh",
				"image_en", "image_zh",
				"trinket_type",
			}),
		},
	).Create(&lesserTrinkets).Error
	if err != nil {
		log.Fatalf("写入小饰品失败: %v", err)
	}
	log.Println("✅ 饰品同步完成")

	// ======================================================
	// 7. 关键词
	// ======================================================
	keywords, err := blizzard.GetKeywords(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取关键词失败: %v", err)
	}
	log.Printf("准备写入关键词数量: %d\n", len(keywords))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh", "text_en", "text_zh",
			}),
		},
	).Create(&keywords).Error
	if err != nil {
		log.Fatalf("写入关键词失败: %v", err)
	}
	log.Println("✅ 关键词同步完成")

	// ======================================================
	// 8. 任务
	// ======================================================
	quests, err := blizzard.GetQuests(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取任务失败: %v", err)
	}
	log.Printf("准备写入任务数量: %d\n", len(quests))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh",
				"text_en", "text_zh",
				"image_en", "image_zh",
			}),
		},
	).Create(&quests).Error
	if err != nil {
		log.Fatalf("写入任务失败: %v", err)
	}
	log.Println("✅ 任务同步完成")

	// ======================================================
	// 9. 奖励
	// ======================================================
	rewards, err := blizzard.GetRewards(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取奖励失败: %v", err)
	}
	log.Printf("准备写入奖励数量: %d\n", len(rewards))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh",
				"text_en", "text_zh",
				"image_en", "image_zh",
			}),
		},
	).Create(&rewards).Error
	if err != nil {
		log.Fatalf("写入奖励失败: %v", err)
	}
	log.Println("✅ 奖励同步完成")

	// ======================================================
	// 10. 法术
	// ======================================================
	spells, err := blizzard.GetSpells(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取法术失败: %v", err)
	}
	log.Printf("准备写入法术数量: %d\n", len(spells))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en", "name_zh",
				"text_en", "text_zh",
				"image_en", "image_zh",
			}),
		},
	).Create(&spells).Error
	if err != nil {
		log.Fatalf("写入法术失败: %v", err)
	}
	log.Println("✅ 法术同步完成")

	// ======================================================
	// 11. Minions（随从）—— 两阶段保存解决外键错误
	// ======================================================

	minions, err := blizzard.GetMinions(tokenResp.AccessToken, nil)
	if err != nil {
		log.Fatalf("获取随从失败: %v", err)
	}
	log.Printf("准备写入随从数量: %d\n", len(minions))

	// 直接 Upsert JSON 模型
	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en",
				"name_zh",
				"text_en",
				"text_zh",
				"attack",
				"health",
				"child_ids",
				"class_ids",
				"image_en",
				"image_zh",
				"is_duo",
				"is_solo",
			}),
		},
	).Create(&minions).Error

	if err != nil {
		log.Fatalf("写入 Minions 失败: %v", err)
	}

	log.Println("✅ 随从 JSON 模式同步完成（无多对多，无外键）")

	// ======================================================
	// 12. 畸变（Anomalies）
	// ======================================================

	anomalies, err := blizzard.GetAnomalies(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取畸变失败: %v", err)
	}
	log.Printf("准备写入畸变数量: %d\n", len(anomalies))

	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en",
				"name_zh",
				"text_en",
				"text_zh",
				"image_en",
				"image_zh",
			}),
		},
	).Create(&anomalies).Error

	if err != nil {
		log.Fatalf("写入畸变失败: %v", err)
	}

	log.Println("✅ 酒馆战旗<畸变 Anomalies> 数据库同步完成")
}
