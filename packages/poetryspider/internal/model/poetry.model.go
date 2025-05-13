package model

type Poetry struct {
	BaseModel
	PoetryId   string `json:"poetry_id" gorm:"column:poetry_id"`     // 诗文ID
	Title      string `json:"title" gorm:"column:title"`             // 名称
	AuthorId   string `json:"author_id" gorm:"column:author_id"`     // 作者
	AuthorName string `json:"author_name" gorm:"column:author_name"` // 作者
	Dynasty    string `json:"dynasty" gorm:"column:dynasty"`         // 朝代
	PoetryType string `json:"poetry_type" gorm:"column:poetry_type"` // 朝代
	Paragraphs string `json:"paragraphs" gorm:"column:paragraphs"`   // 主题
}

func (m *Poetry) TableName() string {
	return "poetry"
}
