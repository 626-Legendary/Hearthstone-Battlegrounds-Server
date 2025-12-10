// models/hero.go
package models

// models/hero.go
type Minions struct {
	ID      uint   
	HSID    int    
	NameEN  string 
	NameZH  string 
	TextEN  string 
	TextZH  string
	Attack int 
	Health int
	ChildIDs int[]
	KeywordIDs int[]
	ImageEN string 
	ImageZH string 
	UpgradeImageEN string // 金卡图片路径
	UpgradeImageZH string
		IsDuo       bool   `gorm:"column:is_duo"`
	IsSolo      bool   `gorm:"column:is_solo"`

}

func (Minions) TableName() string {
	return "minions"
}
