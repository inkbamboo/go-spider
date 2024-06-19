package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Area struct {
	ID           primitive.ObjectID `json:"id" bson:"_id" `
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at" `
	DistrictId   string             `json:"district_id" bson:"district_id" `
	DistrictName string             `json:"district_name" bson:"district_name"`
	AreaName     string             `json:"area_name" bson:"area_name"`
	AreaId       string             `json:"area_id" bson:"area_id"`
}

func (m *Area) TableName() string {
	return "area"
}
func (m *Area) GetBson() (bson.M, error) {
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
