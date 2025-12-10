package models

//任务
type Quests struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	HSID   int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN string `gorm:"column:name_en"`
	NameZH string `gorm:"column:name_zh"`

	TextEN  string
	TextZH  string
	ImageEN string
	ImageZH string
}

func (Quests) TableName() string {
	return "quests"
}
