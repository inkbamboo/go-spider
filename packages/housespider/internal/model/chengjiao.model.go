package model

type ChengJiao struct {
	BaseModel
	HousedelId string  `json:"housedel_id" gorm:"column:housedel_id"` // 房源删除ID
	DistrictId string  `json:"district_id" gorm:"column:district_id"` // 区县名称
	TotalPrice float64 `json:"total_price" gorm:"column:total_price"` // 总价
	UnitPrice  float64 `json:"unit_price" gorm:"column:unit_price"`   // 单价
	DealDate   string  `json:"deal_date" gorm:"column:deal_date"`     // 成交时间
	DealCycle  int64   `json:"deal_cycle" gorm:"column:deal_cycle"`   // 成交周期
	DealPrice  float64 `json:"deal_price" gorm:"column:deal_price"`   // 成交价
}

func (m *ChengJiao) TableName() string {
	return "chengjiao"
}
