package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	houseService     *HouseService
	houseServiceOnce sync.Once
)

func GetHouseService() *HouseService {
	houseServiceOnce.Do(func() {
		houseService = &HouseService{}
	})
	return houseService
}

type HouseService struct {
}

func (s *HouseService) SaveHouse(house *model.House, city string) (err error) {
	tx := ares.Default().GetOrm(city)
	updateInfo := map[string]interface{}{
		"district_id":       house.DistrictId,
		"xiaoqu_name":       house.XiaoquName,
		"house_type":        house.HouseType,
		"house_area":        house.HouseArea,
		"house_orientation": house.HouseOrientation,
		"house_floor":       house.HouseFloor,
	}
	if house.HouseYear != "" {
		updateInfo["house_year"] = house.HouseYear
	}
	if err = tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "housedel_id"}, {Name: "district_id"}},
		DoUpdates: clause.Assignments(updateInfo),
	}).Create(&house).Error; err != nil {
		fmt.Printf("create house error: %v\n", err)
		return
	}
	return
}
