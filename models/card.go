package models

import "time"

// Card 表示一个战棋随从/卡牌
type Card struct {
	ID   uint `json:"id" gorm:"primaryKey"`
	HSID int  `json:"hs_id" gorm:"uniqueIndex"` // 暴雪官方 card id

	Name   string `json:"name"`    // 英文名
	NameZh string `json:"name_zh"` // 中文名

	Tier   int `json:"tier"`   // 酒馆等级
	Attack int `json:"attack"` // 攻击
	Health int `json:"health"` // 生命

	// 种族/类型：比如 "亡灵", "野猪人", "亡灵,野猪人"
	// 这里先用字符串保存，方便一点；以后可以拆表或用 JSON
	Class string `json:"class"` // 文本版，多种族用英文逗号隔开

	// 原始 classId / multiClassIds / childIds 等，用字符串 "1,2,3" 保存
	ClassIDs string `json:"class_ids"` // "1" 或 "1,2"
	ChildIDs string `json:"child_ids"` // 相关卡牌 ID 列表 "123,456"

	Text         string `json:"text"`           // 英文描述
	TextZh       string `json:"text_zh"`        // 中文描述
	FlavorText   string `json:"flavor_text"`    // 彩蛋文本
	FlavorTextZh string `json:"flavor_text_zh"` // 彩蛋文本（中文）

	ImageURL     string `json:"image_url"`      // 普通图
	GoldImageURL string `json:"gold_image_url"` // 金卡图

	// 模式限定
	SolosOnly bool `json:"solos_only"` // 仅单人
	DuosOnly  bool `json:"duos_only"`  // 仅双人

	// 多对多：关键词
	Keywords []Keyword `json:"keywords" gorm:"many2many:card_keywords;"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
