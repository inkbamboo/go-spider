package model

type House struct {
	BaseModel
	HousedelId       string `json:"housedel_id" gorm:"column:housedel_id"`             // 房源删除ID
	DistrictId       string `json:"district_id" gorm:"column:district_id"`             // 区县名称
	HouseArea        string `json:"house_area" gorm:"column:house_area"`               // 房屋面积
	HouseOrientation string `json:"house_orientation" gorm:"column:house_orientation"` // 房屋朝向
	HouseType        string `json:"house_type" gorm:"column:house_type"`               // 房屋类型
	HouseYear        string `json:"house_year" gorm:"column:house_year"`               // 房屋年限
	XiaoquName       string `json:"xiaoqu_name" gorm:"column:xiaoqu_name"`             // 小区名称
	HouseFloor       string `json:"house_floor" gorm:"column:house_floor"`             // 楼层总高度
}

func (m *House) TableName() string {
	return "house"
}
