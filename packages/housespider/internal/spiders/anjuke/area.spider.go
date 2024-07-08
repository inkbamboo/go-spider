package anjuke

import (
	"fmt"
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
		alias: fmt.Sprintf("anjuke_%s", city),
	}
}
func (s *AreaSpider) Start() {
	fmt.Println("Start AreaSpider")
	c := colly.NewCollector()
	//c.OnHTML(".region region-line2", func(e *colly.HTMLElement) {
	c.OnHTML(".region ul", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
		e.ForEach("li", func(_ int, e *colly.HTMLElement) {
			fmt.Println(e.Text)
			href := e.DOM.Find(".anchor anchor-custom").Text()
			fmt.Println(href)
			//if href == "" {
			//	return
			//}
			//fmt.Println(href)
		})
		//if !strings.Contains(href, "sale") {
		//	return
		//}
		//areaId := strings.Split(href, "/")[2]
		//s.parseArea(areaId, e.Text)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(fmt.Sprintf("https://%s.anjuke.com/sale/", s.city))
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
	c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/%s/", s.city, areaId))
}
