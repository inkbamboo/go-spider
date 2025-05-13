package services

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"gorm.io/gorm/clause"
	"sync"
)

var (
	authorService     *AuthorService
	authorServiceOnce sync.Once
)

func GetAuthorService() *AuthorService {
	authorServiceOnce.Do(func() {
		authorService = &AuthorService{}
	})
	return authorService
}

type AuthorService struct {
}

func (s *AuthorService) SaveAuthor(author *model.Author, alias string) (err error) {
	tx := ares.Default().GetOrm(alias)
	if err = tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "author_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":        author.Name,
			"intro":       author.Intro,
			"dynasty":     author.Dynasty,
			"birth_death": author.BirthDeath,
		}),
	}).Create(&author).Error; err != nil {
		fmt.Printf("create author error: %v\n", err)
		return
	}
	return
}
func (s *AuthorService) GetAuthor(alias, authorId string) (author *model.Author, err error) {
	tx := ares.Default().GetOrm(alias)
	err = tx.Model(&model.Author{}).Where("author_id=?", authorId).First(&author).Error
	return
}
func (s *AuthorService) GetAllAuthor(alias string) (authors []*model.Author, err error) {
	tx := ares.Default().GetOrm(alias)
	err = tx.Model(&model.Author{}).Find(&authors).Error
	return
}
