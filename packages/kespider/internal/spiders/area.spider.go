package spiders

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"gorm.io/gorm/clause"
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
	fmt.Println("Start AreaSpider")
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
		areaItem := model.Area{
			AreaId:       areaId,
			AreaName:     areaName,
			DistrictId:   districtId,
			DistrictName: e.Text,
		}
		tx := ares.Default().GetOrm("sjz")
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "district_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"district_name": areaItem.DistrictName,
				"area_id":       areaItem.AreaId,
				"area_name":     areaItem.AreaName,
			}),
		}).Create(&areaItem).Error; err != nil {
			fmt.Printf("create area error: %v\n", err)
			return
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/", areaId))
}
