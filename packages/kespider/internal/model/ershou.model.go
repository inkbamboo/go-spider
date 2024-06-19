package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ErShou struct {
	ID               primitive.ObjectID `json:"id" bson:"_id" `
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at" `
	HousedelId       string             `json:"housedel_id" bson:"housedel_id"`             // 房源删除ID
	AreaName         string             `json:"area_name" bson:"area_name"`                 // 区域名称
	DistrictName     string             `json:"district_name" bson:"district_name"`         // 区县名称
	HouseArea        string             `json:"house_area" bson:"house_area"`               // 房屋面积
	HouseOrientation string             `json:"house_orientation" bson:"house_orientation"` // 房屋朝向
	HouseType        string             `json:"house_type" bson:"house_type"`               // 房屋类型
	HouseYear        string             `json:"house_year" bson:"house_year"`               // 房屋年限
	TotalPrice       string             `json:"total_price" bson:"total_price"`             // 总价
	UnitPrice        string             `json:"unit_price" bson:"unit_price"`               // 单价
	XiaoquName       string             `json:"xiaoqu_name" bson:"xiaoqu_name"`             // 小区名称
	HouseFloor       string             `json:"house_floor" bson:"house_floor"`             // 楼层总高度
}

func (m *ErShou) TableName() string {
	return "ershou"
}
func (m *ErShou) GetBson() (bson.M, error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	m.ID = primitive.NewObjectIDFromTimestamp(time.Now())
	bytes, err := bson.Marshal(m)
	if err != nil {
		return nil, err
	}
	var result bson.M
	err = bson.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
