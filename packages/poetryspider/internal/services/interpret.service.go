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
			"translation": interpret.Translation,
			"intro":       interpret.Intro,
			"annotation":  interpret.Annotation,
		}),
	}).Create(&interpret).Error; err != nil {
		fmt.Printf("create interpret error: %v\n", err)
		return
	}
	return
}
func (s *InterpretService) InterpretExists(poetryId string, alias string) (exists bool) {
	tx := ares.Default().GetOrm(alias)
	var interpret *model.Interpret
	if err := tx.Model(&model.Poetry{}).Where("poetry_id=?", poetryId).First(&interpret).Error; err != nil || interpret == nil {
		return
	}
	return interpret.ID > 0
}
