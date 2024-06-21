package spiders

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/util"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
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

func (s *ErShouSpider) Start() {
	//areas, _ := s.findAllArea()
	//for _, area := range areas {
	//	s.parseOnArea(area)
	//}
	areas, _ := services.GetAreaService().FindAllArea()
	for _, area := range areas {
		if area.DistrictId == "damacun" {
			fmt.Printf("start parse area: %v\n", area.DistrictId)
			s.parseOnArea(area)
		}
	}
}
func (s *ErShouSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector(
		colly.AllowedDomains("sjz.ke.com"), //白名单域名
		colly.AllowURLRevisit(),            //允许对同一 URL 进行多次下载
		colly.Async(true),                  //设置为异步请求
		colly.MaxDepth(2),                  //爬取页面深度,最多为两层
		colly.MaxBodySize(2*1024*1024),     //响应正文最大字节数
		colly.IgnoreRobotsTxt(),            //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	c.OnHTML(".sellListContent", func(e *colly.HTMLElement) {
		s.parseHouseList(area, e)
	})
	c.OnHTML(".page-box div", func(e *colly.HTMLElement) {
		totalPage := gjson.Get(e.Attr("page-data"), "totalPage").Int()
		curPage := gjson.Get(e.Attr("page-data"), "curPage").Int()
		if curPage < totalPage {
			c.UserAgent = ""
			c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/pg%d/", area.DistrictId, curPage+1))
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://sjz.ke.com/ershoufang/%s/", area.DistrictId))
}
func (s *ErShouSpider) parseHouseList(area *model.Area, e *colly.HTMLElement) {
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
		infoStr := util.TrimInfoEmpty(el.DOM.Find(".houseInfo").Text())
		houseItem.HouseArea = util.GetHouseArea(infoStr)
		houseItem.HouseType = util.GetHouseType(infoStr)
		houseItem.HouseFloor = util.GetHouseFloor(infoStr)
		houseItem.HouseOrientation = util.GetHouseOrientation(infoStr)
		houseItem.HouseYear = util.GetHouseYear(infoStr)
		housePrice := &model.HousePrice{
			HousedelId: housedelId,
			Version:    time.Now().Format("2006-01-02"),
			DistrictId: area.DistrictId,
			TotalPrice: util.GetTotalPrice(strings.TrimSpace(el.DOM.Find(".totalPrice").Find("span").Text())),
			UnitPrice:  util.GetUnitPrice(el.DOM.Find(".unitPrice").Find("span").Text()),
		}
		if err := services.GetHouseService().SaveHouse(houseItem); err != nil {
			return
		}
		if err := services.GetHousePriceService().SaveHousePrice(housePrice); err != nil {
			return
		}
	})
}
