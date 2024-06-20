package spider

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
)

type AreaSpider struct{}

var Upsert = true

var (
	areaSpider     *AreaSpider
	areaSpiderOnce sync.Once
)

func GetAreaSpider() *AreaSpider {
	areaSpiderOnce.Do(func() {
		areaSpider = &AreaSpider{}
	})
	return areaSpider
}
func (s *AreaSpider) Start() {
	c := colly.NewCollector()
	c.OnXML("//div[3]/div[1]/dl[2]/dd/div/div/a", func(e *colly.XMLElement) {
		href := e.Attr("href")
		if !strings.Contains(href, "ershoufang") {
			return
		}
		areaId := strings.Split(href, "/")[2]
		s.parseArea(areaId, e.Text)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://sjz.ke.com/ershoufang/")
}
func (s *AreaSpider) parseArea(areaId, areaName string) {
	c := colly.NewCollector()
	c.OnXML("//div[3]/div[1]/dl[2]/dd/div/div[2]/a", func(e *colly.XMLElement) {
		districtId := strings.Split(e.Attr("href"), "/")[2]
		if districtId == "" {
			return
		}
		areaItem := model.Area{}
		filter := bson.D{{"district_id", bson.D{{"$eq", districtId}}}}
		collection := ares.Default().GetMongo("sjz").Collection(areaItem.TableName())
		collection.FindOne(context.TODO(), filter).Decode(&areaItem)
		areaItem.AreaId = areaId
		areaItem.AreaName = areaName
		areaItem.DistrictId = strings.Split(e.Attr("href"), "/")[2]
		areaItem.DistrictName = e.Text
		areaBs, _ := areaItem.ToBson()
		update := bson.M{"$set": areaBs}
		_, _ = collection.UpdateOne(context.TODO(), filter, update, &options.UpdateOptions{Upsert: &Upsert})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/", areaId))
}
