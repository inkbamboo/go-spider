package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
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

func (s *HouseService) SaveHouse(house *model.House) (err error) {
	tx := ares.Default().GetOrm("sjz")
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "housedel_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"district_id":       house.DistrictId,
			"xiaoqu_name":       house.XiaoquName,
			"house_type":        house.HouseType,
			"house_area":        house.HouseArea,
			"house_orientation": house.HouseOrientation,
			"house_year":        house.HouseYear,
			"house_floor":       house.HouseFloor,
		}),
	}).Create(&house).Error; err != nil {
		fmt.Printf("create house error: %v\n", err)
		return
	}
	return
}
