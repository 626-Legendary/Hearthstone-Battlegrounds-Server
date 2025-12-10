// models/hero.go
package models

// models/hero.go
type Keywords struct {
	ID      uint      `gorm:"primaryKey;autoIncrement"`
	HSID    int       `gorm:"column:hs_id;uniqueIndex"`
	NameEN  string    `gorm:"column:name_en"`
	NameZH  string    `gorm:"column:name_zh"`
	TextEN  string    `gorm:"column:text_en"`
	TextZH  string    `gorm:"column:text_zh"`
	ImageEN string    `gorm:"column:image_en"`
	ImageZH string    `gorm:"column:image_zh"`
	Minions []Minions `gorm:"many2many:keyword_minions;"`
}

func (Keywords) TableName() string {
	return "keywords"
}
