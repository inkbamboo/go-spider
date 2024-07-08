package model

type HousePrice struct {
	BaseModel
	HousedelId string  `json:"housedel_id" gorm:"column:housedel_id"` // 房源删除ID
	Version    string  `json:"version" gorm:"column:version"`         // 版本
	DistrictId string  `json:"district_id" gorm:"column:district_id"` // 区县名称
	TotalPrice float64 `json:"total_price" gorm:"column:total_price"` // 总价
	UnitPrice  float64 `json:"unit_price" gorm:"column:unit_price"`   // 单价
}

func (m *HousePrice) TableName() string {
	return "house_price"
}

type XiaoquPrice struct {
	HousedelId   string  `json:"housedel_id" gorm:"column:housedel_id"` // 房源删除ID
	TotalPrice   float64 `json:"total_price" gorm:"column:total_price"` // 总价
	UnitPrice    float64 `json:"unit_price" gorm:"column:unit_price"`   // 单价
	HouseArea    float64 `json:"house_area" gorm:"column:house_area"`   // 房屋面积
	XiaoquName   string  `json:"xiaoqu_name" gorm:"column:xiaoqu_name"` // 小区名称
	DistrictName string  `json:"district_name" gorm:"column:district_name"`
	AreaName     string  `json:"area_name" gorm:"column:area_name"`
}
