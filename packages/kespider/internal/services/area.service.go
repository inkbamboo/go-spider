package services

import (
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
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

func (s *AreaService) FindAllArea() ([]*model.Area, error) {
	tx := ares.Default().GetOrm("sjz")
	var results []*model.Area

	if err := tx.Model(&model.Area{}).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
