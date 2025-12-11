package models

import "gorm.io/datatypes"

// Minions 表示一个随从（酒馆战旗随从）
type Minions struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	HSID   int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN string `gorm:"column:name_en"`
	NameZH string `gorm:"column:name_zh"`
	TextEN string `gorm:"column:text_en"`
	TextZH string `gorm:"column:text_zh"`

	Attack int
	Health int

	// ⭐ childIds 存 JSON
	ChildIDs datatypes.JSON `gorm:"type:json"`

	// ⭐ classIds 也存 JSON
	ClassIDs datatypes.JSON `gorm:"type:json"`

	ImageEN string `gorm:"column:image_en"`
	ImageZH string `gorm:"column:image_zh"`

	IsDuo  bool
	IsSolo bool
}

func (Minions) TableName() string {
	return "minions"
}
