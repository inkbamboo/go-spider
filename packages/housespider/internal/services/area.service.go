package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	areaService     *AreaService
	areaServiceOnce sync.Once
)

func GetAreaService() *AreaService {
	areaServiceOnce.Do(func() {
		areaService = &AreaService{}
	})
	return areaService
}

type AreaService struct {
}

func (s *AreaService) FindAllArea(city string) ([]*model.Area, error) {
	tx := ares.Default().GetOrm(city)
	var results []*model.Area

	if err := tx.Model(&model.Area{}).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
func (s *AreaService) SaveArea(area *model.Area, city string) (err error) {
	tx := ares.Default().GetOrm(city)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "district_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"district_name": area.DistrictName,
			"area_id":       area.AreaId,
			"area_name":     area.AreaName,
		}),
	}).Create(&area).Error; err != nil {
		fmt.Printf("create area error: %v\n", err)
		return
	}
	return
}
