package models

// models/hero.go
type Classes struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	HSID   int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN string `gorm:"column:name_en"`
	NameZH string `gorm:"column:name_zh"`

	Minions []Minion `gorm:"many2many:keyword_minion;"`
}

func (Classes) TableName() string {
	return "classes"
}

/*
classId HSID
1 亡灵
2 野猪人
3 野兽
4 元素
5 机械
7 海盗
8 鱼人
9 恶魔
10 龙
12 中立/英雄/法术
14 娜迦*/
