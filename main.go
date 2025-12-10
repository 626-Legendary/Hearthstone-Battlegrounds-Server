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
	// 读取 .env
	_ = godotenv.Load()

	clientID := os.Getenv("BLIZZARD_CLIENT_ID")
	clientSecret := os.Getenv("BLIZZARD_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("请在 .env 中配置 BLIZZARD_CLIENT_ID 和 BLIZZARD_CLIENT_SECRET")
	}

	// 1. 初始化数据库（建库 + 表）
	database.InitDB()
	db := database.DB

	// 2. 拿 token
	tokenResp, err := blizzard.GetAccessToken(clientID, clientSecret)
	if err != nil {
		log.Fatalf("获取 access_token 失败: %v", err)
	}

	// // 3. 自动翻页获取所有英雄
	// heroes, err := blizzard.GetHeroCards(tokenResp.AccessToken)
	// if err != nil {
	// 	log.Fatalf("获取战棋英雄失败: %v", err)
	// }
	// log.Printf("准备写入英雄数量: %d\n", len(heroes))

	// // 4. 写入英雄：hs_id 存在则更新，不存在则创建
	// err = db.Clauses(
	// 	clause.OnConflict{
	// 		Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
	// 		DoUpdates: clause.AssignmentColumns([]string{
	// 			"name_en",
	// 			"name_zh",
	// 			"armor",
	// 			"hero_power_id",
	// 			"companion_id",
	// 			"image_en",
	// 			"image_zh",
	// 			"is_duo",
	// 			"is_solo",
	// 		}),
	// 	},
	// ).Create(&heroes).Error
	// if err != nil {
	// 	log.Fatalf("写入<英雄数据>到数据库失败: %v", err)
	// }
	// log.Println("✅ 酒馆战旗<英雄数据>数据库同步完成")

	// // 5. 自动翻页获取大饰品
	// greaterTrinkets, err := blizzard.GetGreaterTrinketsCards(tokenResp.AccessToken)
	// if err != nil {
	// 	log.Fatalf("获取大饰品失败: %v", err)
	// }
	// log.Printf("准备写入大饰品数量: %d\n", len(greaterTrinkets))

	// // 6. 写入大饰品：hs_id 存在则更新，不存在则创建
	// err = db.Clauses(
	// 	clause.OnConflict{
	// 		Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
	// 		DoUpdates: clause.AssignmentColumns([]string{
	// 			"name_en",
	// 			"name_zh",
	// 			"mana_cost",
	// 			"text_en",
	// 			"text_zh",
	// 			"image_en",
	// 			"image_zh",
	// 			"trinket_type",
	// 		}),
	// 	},
	// ).Create(&greaterTrinkets).Error
	// if err != nil {
	// 	log.Fatalf("写入<大饰品数据>到数据库失败: %v", err)
	// }

	// // 7. 自动翻页获取小饰品
	// lesserTrinkets, err := blizzard.GetLesserTrinketsCards(tokenResp.AccessToken)
	// if err != nil {
	// 	log.Fatalf("获取小饰品失败: %v", err)
	// }
	// log.Printf("准备写入小饰品数量: %d\n", len(lesserTrinkets))

	// // 8. 写入小饰品：hs_id 存在则更新，不存在则创建
	// err = db.Clauses(
	// 	clause.OnConflict{
	// 		Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
	// 		DoUpdates: clause.AssignmentColumns([]string{
	// 			"name_en",
	// 			"name_zh",
	// 			"mana_cost",
	// 			"text_en",
	// 			"text_zh",
	// 			"image_en",
	// 			"image_zh",
	// 			"trinket_type",
	// 		}),
	// 	},
	// ).Create(&lesserTrinkets).Error
	// if err != nil {
	// 	log.Fatalf("写入<小饰品数据>到数据库失败: %v", err)
	// }
	// log.Println("✅ 酒馆战旗<饰品数据>数据库同步完成")

	// 9. 获取关键词（不需要翻页）
	keywords, err := blizzard.GetKeywords(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取关键词失败: %v", err)
	}
	log.Printf("准备写入关键词数量: %d\n", len(keywords))

	// 10. 写入关键词：hs_id 存在则更新，不存在则创建
	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en",
				"name_zh",
				"text_en",
				"text_zh",
			}),
		},
	).Create(&keywords).Error
	if err != nil {
		log.Fatalf("写入<关键词数据>到数据库失败: %v", err)
	}
	log.Println("✅ 酒馆战旗<关键词数据>数据库同步完成")

	//11. quest
	quests, err := blizzard.GetQuests(tokenResp.AccessToken)
	if err != nil {
		log.Fatalf("获取任务失败: %v", err)
	}
	log.Printf("准备写入任务数量: %d\n", len(quests))
	// 12 databse quest
	err = db.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "hs_id"}}, // 根据 hs_id 判断冲突
			DoUpdates: clause.AssignmentColumns([]string{
				"name_en",
				"name_zh",
				"text_en",
				"text_zh",
				"image_en",
				"image_zh",
			}),
		},
	).Create(&quests).Error
	if err != nil {
		log.Fatalf("写入<任务数据>到数据库失败: %v", err)
	}
	log.Println("✅ 酒馆战旗<任务数据>数据库同步完成")

}
