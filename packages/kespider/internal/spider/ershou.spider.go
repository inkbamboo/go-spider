package spider

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/util"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"gorm.io/gorm/clause"
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
	tx := ares.Default().GetOrm("sjz")
	var results []*model.Area

	if err := tx.Model(&model.Area{}).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (s *ErShouSpider) Start() {
	//areas, _ := s.findAllArea()
	//for _, area := range areas {
	//	s.parseOnArea(area)
	//}
	areas, _ := s.findAllArea()
	for _, area := range areas {
		if area.DistrictId == "damacun" {
			s.parseHouseList(area, 1)
		}
	}
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
			housedelId, _ := lo.Last(strings.Split(href, "/"))
			housedelId = strings.TrimSuffix(housedelId, ".html")
			if housedelId == "" {
				return
			}
			houseItem := &model.House{
				HousedelId: housedelId,
				DistrictId: area.DistrictId,
				XiaoquName: el.DOM.Find(".positionInfo").Find("a").Text(),
			}
			houseItem.HouseArea, houseItem.HouseType, houseItem.HouseOrientation, houseItem.HouseYear, houseItem.HouseFloor = util.ParseHouseDetail(el.DOM.Find(".houseInfo").Text())
			tx := ares.Default().GetOrm("sjz")
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "housedel_id"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"district_id":       houseItem.DistrictId,
					"xiaoqu_name":       houseItem.XiaoquName,
					"house_type":        houseItem.HouseType,
					"house_area":        houseItem.HouseArea,
					"house_orientation": houseItem.HouseOrientation,
					"house_year":        houseItem.HouseYear,
					"house_floor":       houseItem.HouseFloor,
				}),
			}).Create(&houseItem).Error; err != nil {
				fmt.Printf("create area error: %v\n", err)
				return
			}
			totalPrice := strings.TrimSpace(el.DOM.Find(".totalPrice").Find("span").Text())
			unitPrice := el.DOM.Find(".unitPrice").Find("span").Text()
			unitPrice = strings.ReplaceAll(unitPrice, "元/平", "")
			unitPrice = strings.TrimSpace(strings.ReplaceAll(unitPrice, ",", ""))
			housePrice := &model.HousePrice{
				HousedelId: housedelId,
				Version:    time.Now().Format("2006-01-02"),
				DistrictId: area.DistrictId,
				TotalPrice: cast.ToFloat64(totalPrice),
				UnitPrice:  cast.ToFloat64(unitPrice),
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "housedel_id"}, {Name: "version"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"district_id": housePrice.DistrictId,
					"total_price": housePrice.TotalPrice,
					"unit_price":  housePrice.UnitPrice,
				}),
			}).Create(&housePrice).Error; err != nil {
				fmt.Printf("create priceInfo error: %v\n", err)
				return
			}
		})
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/pg%d/", area.DistrictId, page))
}
