package models

//任务
type Quests struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	HSID   int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN string `gorm:"column:name_en"`
	NameZH string `gorm:"column:name_zh"`

	TextEN  string `gorm:"column:text_en"`
	TextZH  string `gorm:"column:text_zh"`
	ImageEN string `gorm:"column:image_en"`
	ImageZH string `gorm:"column:image_zh"`
}

func (Quests) TableName() string {
	return "quests"
}

/*
https://us.api.blizzard.com/hearthstone/cards?type=spell&gameMode=battlegrounds&sort=quest
*/
