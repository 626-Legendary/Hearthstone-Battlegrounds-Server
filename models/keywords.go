package models

// Keyword 表示一个战棋关键词（嘲讽、圣盾、亡语……）
type Keyword struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	HSID   int    `json:"hs_id" gorm:"uniqueIndex"` // 暴雪 keywordId
	Name   string `json:"name"`                     // 英文名
	NameZh string `json:"name_zh"`                  // 中文名（如果你以后需要的话）

	// 反向关联，可选
	Cards []Card `json:"-" gorm:"many2many:card_keywords;"`
}
