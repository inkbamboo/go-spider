package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	interpretService     *InterpretService
	interpretServiceOnce sync.Once
)

func GetInterpretService() *InterpretService {
	interpretServiceOnce.Do(func() {
		interpretService = &InterpretService{}
	})
	return interpretService
}

type InterpretService struct {
}

func (s *InterpretService) SaveInterpret(interpret *model.Interpret, alias string) (err error) {
	tx := ares.Default().GetOrm(alias)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "poetry_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"translation":      interpret.Translation,
			"evaluate":         interpret.Evaluate,
			"explanatory_note": interpret.ExplanatoryNote,
		}),
	}).Create(&interpret).Error; err != nil {
		fmt.Printf("create interpret error: %v\n", err)
		return
	}
	return
}
