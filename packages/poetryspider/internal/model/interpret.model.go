package model

type Interpret struct {
	BaseModel
	PoetryId        string `json:"poetry_id" gorm:"column:poetry_id"`               // 诗文ID
	Translation     string `json:"translation" gorm:"column:translation"`           // 译文
	ExplanatoryNote string `json:"explanatory_note" gorm:"column:explanatory_note"` // 注释
	Evaluate        string `json:"evaluate" gorm:"column:evaluate"`                 // 评价
}

func (m *Interpret) TableName() string {
	return "interpret"
}
