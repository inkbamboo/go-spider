package spiders

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/util"
	"github.com/samber/lo"
	"github.com/spf13/cast"
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
		//colly.Debugger(&debug.LogDebugger{}), // 开启debug
		colly.MaxDepth(2),              //爬取页面深度,最多为两层
		colly.MaxBodySize(2*1024*1024), //响应正文最大字节数
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "),
		colly.IgnoreRobotsTxt(), //忽略目标机器中的`robots.txt`声明
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
		houseItem.HouseArea, houseItem.HouseType, houseItem.HouseOrientation, houseItem.HouseYear, houseItem.HouseFloor = util.ParseHouseDetail(el.DOM.Find(".houseInfo").Text())

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
		if err := services.GetHouseService().SaveHouse(houseItem); err != nil {
			return
		}
		if err := services.GetHousePriceService().SaveHousePrice(housePrice); err != nil {
			return
		}
	})
}
