package scripts

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders"
	"time"
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
func getOneBatchPoetry(startId int64, batchSize int) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Poetry
	_ = stageTx.Model(&model.Poetry{}).Where("id>?", startId).Order("id asc").Limit(batchSize).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		time.Sleep(3 * time.Millisecond)
		if services.GetPoetryService().PoetryExists(item.PoetryId, "zhsc_poetry") {
			continue
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

	return len(list) == batchSize, endId
}
func getOneBatchAuthor(startId int64, batchSize int) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Author
	_ = stageTx.Model(&model.Author{}).Where("id>?", startId).Order("id asc").Limit(batchSize).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		time.Sleep(3 * time.Millisecond)
		if services.GetAuthorService().AuthorExists(item.AuthorId, "zhsc_poetry") {
			continue
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

	return len(list) == batchSize, endId
}
func getOneBatchInterpret(startId int64, batchSize int) (hasNext bool, endId int64) {
	stageTx := ares.Default().GetOrm("stage_zhsc_poetry")
	var list []*model.Interpret
	_ = stageTx.Model(&model.Interpret{}).Where("id>?", startId).Order("id asc").Limit(batchSize).Find(&list).Error
	for _, item := range list {
		if item.ID > endId {
			endId = item.ID
		}
		time.Sleep(3 * time.Millisecond)
		if services.GetInterpretService().InterpretExists(item.PoetryId, "zhsc_poetry") {
			continue
		}
		_ = services.GetInterpretService().SaveInterpret(&model.Interpret{
			PoetryId:    item.PoetryId,
			Translation: item.Translation,
			Annotation:  item.Annotation,
			Intro:       item.Intro,
		}, "zhsc_poetry")
		fmt.Printf("************ %+v  %+v \n", item.ID, item.PoetryId)

	}

	return len(list) == batchSize, endId
}
func RunTest() {
	//var startId int64
	//batchSize := 2000
	//for {
	//	var hasNext bool
	//	//if hasNext, startId = getOneBatchPoetry(startId, batchSize); !hasNext {
	//	//	break
	//	//}
	//	//if hasNext, startId = getOneBatchAuthor(startId, batchSize); !hasNext {
	//	//	break
	//	//}
	//	if hasNext, startId = getOneBatchInterpret(startId, batchSize); !hasNext {
	//		break
	//	}
	//}
	select {}
}
