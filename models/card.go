package models

import "time"

// Card 表示一个战棋随从/卡牌 不需要写入数据库
type Card struct {
	HSID int `json:"hs_id" ` // 暴雪官方 card id

	ClassID       int   `json:"classId"`
	MultiClassIds []int `json:"multiClassIds"` // 多职业卡牌的其他职业 ID 列表

	ManaCost int               `json:"manacost"` // 花费
	Name     map[string]string `json:"name"`     // 卡牌名称，key 为语言代号

	Tier   int `json:"tier"`   // 酒馆等级
	Attack int `json:"attack"` // 攻击
	Health int `json:"health"` // 生命

	// 文本版种族/类型，多种族用英文逗号隔开
	Class string `json:"class"`

	// 关联表：Classes（多对多）
	Classes []Classes `json:"classes" gorm:"many2many:card_classes;"`

	// 相关卡牌 ID 列表，CSV 格式
	ChildIDs string `json:"child_ids"`

	TextEN       string `json:"text_en"`        // 英文描述
	TextZH       string `json:"text_zh"`        // 中文描述
	FlavorTextEN string `json:"flavor_text_en"` // 彩蛋文本
	FlavorTextZH string `json:"flavor_text_zh"` // 彩蛋文本（中文）

	ImageEN     string `json:"image_en"`      // 普通图
	ImageZH     string `json:"image_zh"`      // 普通图
	GoldImageEN string `json:"gold_image_en"` // 金卡图
	GoldImageZH string `json:"gold_image_zh"` // 金卡图

	// 模式限定
	SolosOnly bool `json:"solos_only"` // 仅单人
	DuosOnly  bool `json:"duos_only"`  // 仅双人

	// 多对多：关键词

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Card) TableName() string {
	return "cards"
}
