package model

type Interpret struct {
	BaseModel
	PoetryId    string `json:"poetry_id" gorm:"column:poetry_id"`     // 诗文ID
	Translation string `json:"translation" gorm:"column:translation"` // 译文
	Annotation  string `json:"annotation" gorm:"column:annotation"`   // 注释
	Intro       string `json:"intro" gorm:"column:intro"`             // 评价
}

func (m *Interpret) TableName() string {
	return "interpret"
}
