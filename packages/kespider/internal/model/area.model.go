package model

type Area struct {
	BaseModel
	DistrictId   string `json:"district_id" gorm:"column:district_id" `
	DistrictName string `json:"district_name" gorm:"column:district_name"`
	AreaName     string `json:"area_name" gorm:"column:area_name"`
	AreaId       string `json:"area_id" gorm:"column:area_id"`
}

func (m *Area) TableName() string {
	return "area"
}
