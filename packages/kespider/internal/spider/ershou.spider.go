package spider

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/util"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"sync"
	"time"
)

type ErShouSpider struct{}

var (
	erShouSpider     *ErShouSpider
	erShouSpiderOnce sync.Once
)

func GetErShouSpider() *ErShouSpider {
	erShouSpiderOnce.Do(func() {
		erShouSpider = &ErShouSpider{}
	})
	return erShouSpider
}
func (s *ErShouSpider) findAllArea() ([]*model.Area, error) {
	cursor, err := ares.Default().GetMongo("sjz").Collection("area").Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, err
	}
	// Unpacks the cursor into a slice
	var results []*model.Area
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results, nil
}

func (s *ErShouSpider) Start() {
	//areas, _ := s.findAllArea()
	//for _, area := range areas {
	//	s.parseOnArea(area)
	//}
	s.parseHouseList(&model.Area{DistrictId: "damacun"}, 1)
}
func (s *ErShouSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector()
	c.OnHTML(".page-box div", func(e *colly.HTMLElement) {
		totalPage := gjson.Get(e.Attr("page-data"), "totalPage").Int()
		fmt.Printf("Found %v\n", totalPage)
		for i := int64(1); i <= totalPage; i++ {
			s.parseHouseList(area, i)
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/", area.DistrictId))
}
func (s *ErShouSpider) parseHouseList(area *model.Area, page int64) {
	c := colly.NewCollector()
	c.OnHTML(".sellListContent", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, el *colly.HTMLElement) {
			href, _ := el.DOM.Find("a").Attr("href")
			if !strings.HasSuffix(href, ".html") {
				return
			}
			houseItem := &model.ErShouFang{}
			houseItem.AreaName = area.AreaName
			houseItem.DistrictName = area.DistrictName
			houseItem.HousedelId, _ = lo.Last(strings.Split(href, "/"))
			houseItem.HousedelId = strings.TrimSuffix(houseItem.HousedelId, ".html")
			houseItem.XiaoquName = el.DOM.Find(".positionInfo").Find("a").Text()
			priceInfo := &model.PriceInfo{}
			priceInfo.TotalPrice = el.DOM.Find(".totalPrice").Find("span").Text()
			priceInfo.UnitPrice = el.DOM.Find(".unitPrice").Find("span").Text()
			priceInfo.UnitPrice = strings.ReplaceAll(priceInfo.UnitPrice, "元/平", "")
			priceInfo.UnitPrice = strings.ReplaceAll(priceInfo.UnitPrice, ",", "")
			priceInfo.DateStr = time.Now().Format("2006-01-02")
			houseType, houseArea, houseOrientation, houseYear, houseFloor := util.ParseHouseDetail(el.DOM.Find(".houseInfo").Text())
			houseItem.HouseType = houseType
			houseItem.HouseArea = houseArea
			houseItem.HouseOrientation = houseOrientation
			houseItem.HouseYear = houseYear
			houseItem.HouseFloor = houseFloor
			filter := bson.D{{"housedel_id", bson.D{{"$eq", houseItem.HousedelId}}}}
			houseBs, _ := houseItem.GetBson()
			update := bson.M{"$set": houseBs}
			ares.Default().GetMongo("sjz").Collection(houseItem.TableName()).UpdateOne(context.TODO(), filter, update, &options.UpdateOptions{Upsert: &Upsert})
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/pg%d/", area.DistrictId, page))
}
