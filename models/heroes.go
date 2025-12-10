package models

// Heroes 英雄

type Heroes struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	HSID        int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN      string `gorm:"column:name_en"`
	NameZH      string `gorm:"column:name_zh"`
	Armor       int    `gorm:"column:armor"`
	HeroPowerID int    `gorm:"column:hero_power_id"`
	CompanionID int    `gorm:"column:companion_id"`
	ImageEN     string `gorm:"column:image_en"`
	ImageZH     string `gorm:"column:image_zh"`
	IsDuo       bool   `gorm:"column:is_duo"`
	IsSolo      bool   `gorm:"column:is_solo"`
}

func (Heroes) TableName() string {
	return "heroes"
}
