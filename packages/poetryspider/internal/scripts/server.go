package scripts

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders"
)

func RunPoetrySpider(platform, spider string) {
	fmt.Println("RunPoetrySpider")
	if sp, err := spiders.NewInstance(platform, spider); err != nil {
		fmt.Println("RunPoetrySpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
}
func getOneBatchPoetry(startId int64) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Poetry
	_ = stageTx.Model(&model.Poetry{}).Where("id>?", startId).Order("id asc").Limit(1000).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		_ = services.GetPoetryService().SavePoetry(&model.Poetry{
			PoetryId:   item.PoetryId,
			Title:      item.Title,
			AuthorId:   item.AuthorId,
			AuthorName: item.AuthorName,
			Dynasty:    item.Dynasty,
			PoetryType: item.PoetryType,
			Paragraphs: item.Paragraphs,
		}, "zhsc_poetry")

		fmt.Printf("************ %+v %+v  %+v\n", item.ID, item.PoetryId, item.Title)
	}

	return len(list) == 1000, endId
}
func getOneBatchAuthor(startId int64) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Author
	_ = stageTx.Model(&model.Author{}).Where("id>?", startId).Order("id asc").Limit(1000).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		_ = services.GetAuthorService().SaveAuthor(&model.Author{
			AuthorId:   item.AuthorId,
			AuthorName: item.AuthorName,
			Dynasty:    item.Dynasty,
			BirthDeath: item.BirthDeath,
			Intro:      item.Intro,
		}, "zhsc_poetry")
		fmt.Printf("************ %+v %+v  %+v\n", item.ID, item.AuthorId, item.AuthorName)
	}

	return len(list) == 1000, endId
}
func getOneBatchInterpret(startId int64) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Interpret
	_ = stageTx.Model(&model.Interpret{}).Where("id>?", startId).Order("id asc").Limit(1000).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		_ = services.GetInterpretService().SaveInterpret(&model.Interpret{
			PoetryId:    item.PoetryId,
			Translation: item.Translation,
			Annotation:  item.Annotation,
			Intro:       item.Intro,
		}, "zhsc_poetry")
		fmt.Printf("************ %+v  %+v \n", item.ID, item.PoetryId)
	}

	return len(list) == 1000, endId
}
func RunTest() {
	var startId int64
	for {
		var hasNext bool
		//if hasNext, startId = getOneBatchPoetry(startId); !hasNext {
		//	break
		//}
		//if hasNext, startId = getOneBatchAuthor(startId); !hasNext {
		//	break
		//}
		if hasNext, startId = getOneBatchInterpret(startId); !hasNext {
			break
		}
	}
	//spider := zhsc.NewPoetrySpider()
	//spider.Test()
	select {}
}
