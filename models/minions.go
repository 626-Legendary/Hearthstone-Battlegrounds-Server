package models

// Minions 表示一个随从（酒馆战旗随从）
type Minions struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	HSID           int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN         string `gorm:"column:name_en"`
	NameZH         string `gorm:"column:name_zh"`
	TextEN         string `gorm:"column:text_en"`
	TextZH         string `gorm:"column:text_zh"`
	Attack         int
	Health         int
	ChildIDs       string `gorm:"column:child_ids"` // CSV
	ImageEN        string `gorm:"column:image_en"`
	ImageZH        string `gorm:"column:image_zh"`
	UpgradeImageEN string `gorm:"column:upgrade_image_en"` // 金卡图片路径
	UpgradeImageZH string `gorm:"column:upgrade_image_zh"`
	IsDuo          bool   `gorm:"column:is_duo"`
	IsSolo         bool   `gorm:"column:is_solo"`

	Classes  []Classes  `gorm:"many2many:class_minions;"`
	Keywords []Keywords `gorm:"many2many:keyword_minions;"`
}

func (Minions) TableName() string {
	return "minions"
}
