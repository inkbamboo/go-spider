package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	housePriceService     *HousePriceService
	housePriceServiceOnce sync.Once
)

func GetHousePriceService() *HousePriceService {
	housePriceServiceOnce.Do(func() {
		housePriceService = &HousePriceService{}
	})
	return housePriceService
}

type HousePriceService struct {
}

func (s *HousePriceService) SaveHousePrice(housePrice *model.HousePrice) (err error) {
	tx := ares.Default().GetOrm("sjz")
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "housedel_id"}, {Name: "version"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"district_id": housePrice.DistrictId,
			"total_price": housePrice.TotalPrice,
			"unit_price":  housePrice.UnitPrice,
		}),
	}).Create(&housePrice).Error; err != nil {
		fmt.Printf("create housePrice error: %v\n", err)
		return
	}
	return
}
