package model

type HousePrice struct {
	BaseModel
	HousedelId string  `json:"housedel_id" gorm:"column:housedel_id"` // 房源删除ID
	Version    string  `json:"version" gorm:"column:version"`         // 版本
	DistrictId string  `json:"district_id" gorm:"column:district_id"` // 区县名称
	TotalPrice float64 `json:"total_price" bson:"total_price"`        // 总价
	UnitPrice  float64 `json:"unit_price" bson:"unit_price"`          // 单价
}

func (m *HousePrice) TableName() string {
	return "house_price"
}
