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
			"author_name": author.AuthorName,
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
func (s *AuthorService) GetAuthor(authorId, alias string) (author *model.Author, err error) {
	tx := ares.Default().GetOrm(alias)
	err = tx.Model(&model.Author{}).Where("author_id=?", authorId).First(&author).Error
	return
}
func (s *AuthorService) AuthorExists(authorId string, alias string) (exists bool) {
	tx := ares.Default().GetOrm(alias)
	var author *model.Author
	if err := tx.Model(&model.Author{}).Where("author_id=?", authorId).First(&author).Error; err != nil || author == nil {
		return
	}
	return author.ID > 0
}
func (s *AuthorService) GetAllAuthor(alias string) (authors []*model.Author, err error) {
	tx := ares.Default().GetOrm(alias)
	err = tx.Model(&model.Author{}).Find(&authors).Error
	return
}
