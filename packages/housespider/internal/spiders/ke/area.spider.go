package ke

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/services"
	"strings"
)

type AreaSpider struct {
	city  string
	alias string
}

func NewAreaSpider(city string) *AreaSpider {
	return &AreaSpider{
		city:  city,
		alias: fmt.Sprintf("ke_%s", city),
	}
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

	c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/", s.city))
}
func (s *AreaSpider) parseArea(areaId, areaName string) {
	c := colly.NewCollector()
	c.OnXML("//div[3]/div[1]/dl[2]/dd/div/div[2]/a", func(e *colly.XMLElement) {
		districtId := strings.Split(e.Attr("href"), "/")[2]
		if districtId == "" {
			return
		}
		areaItem := &model.Area{
			AreaId:       areaId,
			AreaName:     areaName,
			DistrictId:   districtId,
			DistrictName: e.Text,
		}
		services.GetAreaService().SaveArea(areaItem, s.alias)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.UserAgent = browser.Random()
	c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/%s/", s.city, areaId))
}
