package models

// models/hero.go
type Classes struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	HSID   int    `gorm:"column:hs_id;uniqueIndex"`
	NameEN string `gorm:"column:name_en"`
	NameZH string `gorm:"column:name_zh"`

	Minions []Minions `gorm:"many2many:class_minions;"`
}

func (Classes) TableName() string {
	return "classes"
}

/*
classId HSID
1 亡灵 undead
2 野猪人 quilboar
3 野兽  beast
4 元素  elemental
5 机械  mech
7 海盗  pirate
8 鱼人  murloc
9 恶魔  demon
10 龙   dragon
12 中立 all
14 娜迦 naga
*/

// GetClasses 返回预置的职业/种族列表，用于初始化 classes 表
func GetClasses() []Classes {
	return []Classes{
		{HSID: 1, NameEN: "undead", NameZH: "亡灵"},
		{HSID: 2, NameEN: "quilboar", NameZH: "野猪人"},
		{HSID: 3, NameEN: "beast", NameZH: "野兽"},
		{HSID: 4, NameEN: "elemental", NameZH: "元素"},
		{HSID: 5, NameEN: "mech", NameZH: "机械"},
		{HSID: 7, NameEN: "pirate", NameZH: "海盗"},
		{HSID: 8, NameEN: "murloc", NameZH: "鱼人"},
		{HSID: 9, NameEN: "demon", NameZH: "恶魔"},
		{HSID: 10, NameEN: "dragon", NameZH: "龙"},
		{HSID: 12, NameEN: "all", NameZH: "中立"},
		{HSID: 14, NameEN: "naga", NameZH: "娜迦"},
	}
}
