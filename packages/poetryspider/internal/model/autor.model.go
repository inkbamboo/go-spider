package model

type Author struct {
	BaseModel
	AuthorId   string `json:"author_id" gorm:"column:author_id"`     // 诗文ID
	Name       string `json:"name" gorm:"column:name"`               // 名称
	Dynasty    string `json:"dynasty" gorm:"column:dynasty"`         // 朝代
	BirthDeath string `json:"birth_death" gorm:"column:birth_death"` // 朝代
	Intro      string `json:"intro" gorm:"column:intro"`             // 评价
}

func (m *Author) TableName() string {
	return "author"
}
