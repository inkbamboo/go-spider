package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	chengJiaoService     *ChengJiaoService
	chengJiaoServiceOnce sync.Once
)

func GetChengJiaoService() *ChengJiaoService {
	chengJiaoServiceOnce.Do(func() {
		chengJiaoService = &ChengJiaoService{}
	})
	return chengJiaoService
}

type ChengJiaoService struct {
}

func (s *ChengJiaoService) SaveChengJiao(chengJiao *model.ChengJiao, alias string) (err error) {
	tx := ares.Default().GetOrm(alias)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "housedel_id"}, {Name: "district_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"district_id": chengJiao.DistrictId,
			"total_price": chengJiao.TotalPrice,
			"unit_price":  chengJiao.UnitPrice,
			"deal_date":   chengJiao.DealDate,
			"deal_price":  chengJiao.DealPrice,
			"deal_cycle":  chengJiao.DealCycle,
		}),
	}).Create(&chengJiao).Error; err != nil {
		fmt.Printf("create chengJiao error: %v\n", err)
		return
	}
	return
}
