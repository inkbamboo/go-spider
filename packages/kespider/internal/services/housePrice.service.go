package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/samber/lo"
	"gorm.io/gorm/clause"
	"strings"
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

func (s *HousePriceService) SaveHousePrice(housePrice *model.HousePrice, city string) (err error) {
	tx := ares.Default().GetOrm(city)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "housedel_id"}, {Name: "district_id"}, {Name: "version"}},
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
func (s *HousePriceService) GetChangeHouse(versions []string) string {
	tx := ares.Default().GetOrm("sjz")
	var hosePrices []*model.HousePrice
	_ = tx.Model(&model.HousePrice{}).Where("version in(?)", versions).Find(&hosePrices).Error
	houseInfos := lo.GroupBy(hosePrices, func(item *model.HousePrice) string {
		return item.HousedelId + item.DistrictId
	})
	changeHouse := fmt.Sprintf("'housedel_id,%s", strings.Join(versions, ","))
	changeHouse += fmt.Sprintf("'\n")
	for _, houseInfo := range houseInfos {
		housePrice := map[string]float64{}
		for _, item := range houseInfo {
			housePrice[item.Version] = item.TotalPrice
		}
		line := fmt.Sprintf("'%s", houseInfo[0].HousedelId)
		var curlPrice float64
		var isChange bool
		for _, version := range versions {
			line += fmt.Sprintf(",%v", housePrice[version])
			if housePrice[version] == 0 {
				continue
			}
			if curlPrice == 0 {
				curlPrice = housePrice[version]
			}
			if !isChange && curlPrice > housePrice[version] {
				isChange = true
			}
		}
		line += fmt.Sprintf("'\n")
		if isChange {
			changeHouse += line
		}
	}
	return changeHouse
}
