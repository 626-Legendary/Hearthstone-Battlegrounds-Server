// 法术
package models

type Spells struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	HSID     int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN   string `gorm:"column:name_en"`
	NameZH   string `gorm:"column:name_zh"`
	ManaCost int    `gorm:"column:mana_cost"`
	TextEN   string `gorm:"column:text_en"`
	TextZH   string `gorm:"column:text_zh"`
	ImageEN  string `gorm:"column:image_en"`
	ImageZH  string `gorm:"column:image_zh"`
}

func (Spells) TableName() string {
	return "spells"
}
