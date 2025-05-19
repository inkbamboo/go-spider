package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	poetryService     *PoetryService
	poetryServiceOnce sync.Once
)

func GetPoetryService() *PoetryService {
	poetryServiceOnce.Do(func() {
		poetryService = &PoetryService{}
	})
	return poetryService
}

type PoetryService struct {
}

func (s *PoetryService) SavePoetry(poetry *model.Poetry, alias string) (err error) {
	tx := ares.Default().GetOrm(alias)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "poetry_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":       poetry.Title,
			"author_name": poetry.AuthorName,
			"author_id":   poetry.AuthorId,
			"dynasty":     poetry.Dynasty,
			"paragraphs":  poetry.Paragraphs,
		}),
	}).Create(&poetry).Error; err != nil {
		fmt.Printf("create poetry error: %v\n", err)
		return
	}
	return
}
func (s *PoetryService) PoetryExists(poetryId string, alias string) (exists bool) {
	tx := ares.Default().GetOrm(alias)
	var poetry *model.Poetry
	if err := tx.Model(&model.Poetry{}).Where("poetry_id=?", poetryId).First(&poetry).Error; err != nil || poetry == nil {
		return
	}
	return poetry.ID > 0
}
