package models

type Hero struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	HSID        int    `gorm:"column:hs_id"` // Blizzard 卡牌 ID
	NameEN      string `gorm:"column:name_en"`
	NameZH      string `gorm:"column:name_zh"`
	Armor       int
	HeroPowerID int    `gorm:"column:hero_power_id"`
	CompanionID int    `gorm:"column:companion_id"`
	ImageEN     string `gorm:"column:image_en"`
	ImageZH     string `gorm:"column:image_zh"`
	IsDuo       bool   `gorm:"column:is_duo"`
	IsSolo      bool   `gorm:"column:is_solo"`
}

// 设置表名（可选）
func (Hero) TableName() string {
	return "heroes"
}
