package ke

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/util"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type ErShouSpider struct {
	city  string
	alias string
}

func NewErShouSpider(city string) *ErShouSpider {
	return &ErShouSpider{
		city:  city,
		alias: fmt.Sprintf("ke_%s", city),
	}
}

func (s *ErShouSpider) Start() {
	areas, _ := services.GetAreaService().FindAllArea(s.alias)
	for _, area := range areas {
		s.parseOnArea(area)
		time.Sleep(10 * time.Second)
	}
	//areas, _ := services.GetAreaService().FindAllArea(s.alias)
	//for _, area := range areas {
	//	if area.DistrictId == "damacun" {
	//		fmt.Printf("start parse area: %v\n", area.DistrictId)
	//		s.parseOnArea(area)
	//	}
	//}
}
func (s *ErShouSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector(
		colly.AllowedDomains(fmt.Sprintf("%s.ke.com", s.city)), //白名单域名
		colly.AllowURLRevisit(),                                //允许对同一 URL 进行多次下载
		colly.Async(true),                                      //设置为异步请求
		colly.MaxDepth(2),                                      //爬取页面深度,最多为两层
		colly.MaxBodySize(2*1024*1024),                         //响应正文最大字节数
		colly.IgnoreRobotsTxt(),                                //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		RandomDelay: 10 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	c.OnHTML(".leftContent", func(e *colly.HTMLElement) {
		s.parseHouseList(area, e.DOM.Find(".sellListContent").Find("li"))
		pageData, _ := e.DOM.Find(".page-box div").Attr("page-data")
		totalPage := gjson.Get(pageData, "totalPage").Int()
		curPage := gjson.Get(pageData, "curPage").Int()
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/%s/pg%d/", s.city, area.DistrictId, curPage+1))
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/%s/", s.city, area.DistrictId))
}
func (s *ErShouSpider) parseHouseList(area *model.Area, e *goquery.Selection) {
	e.Each(func(_ int, el *goquery.Selection) {
		href, _ := el.Find("a").Attr("href")
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
			XiaoquName: el.Find(".positionInfo").Find("a").Text(),
		}
		infoStr := util.TrimInfoEmpty(el.Find(".houseInfo").Text())
		houseItem.HouseArea = util.GetHouseArea(infoStr)
		houseItem.HouseType = util.GetHouseType(infoStr)
		houseItem.HouseFloor = util.GetHouseFloor(infoStr)
		houseItem.HouseOrientation = util.GetHouseOrientation(infoStr)
		houseItem.HouseYear = util.GetHouseYear(infoStr)
		housePrice := &model.HousePrice{
			HousedelId: housedelId,
			Version:    time.Now().Format("2006-01-02"),
			DistrictId: area.DistrictId,
			TotalPrice: util.GetTotalPrice(strings.TrimSpace(el.Find(".totalPrice.totalPrice2").Find("span").Text())),
			UnitPrice:  util.GetUnitPrice(el.Find(".unitPrice").Find("span").Text()),
		}
		if err := services.GetHouseService().SaveHouse(houseItem, s.alias); err != nil {
			return
		}
		if err := services.GetHousePriceService().SaveHousePrice(housePrice, s.alias); err != nil {
			return
		}
	})
}
