package spiders

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/ares"
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

type ChengJiaoSpider struct{}

var (
	chengJiaoSpider     *ChengJiaoSpider
	chengJiaoSpiderOnce sync.Once
)

func GetChengJiaoSpider() *ChengJiaoSpider {
	chengJiaoSpiderOnce.Do(func() {
		chengJiaoSpider = &ChengJiaoSpider{}
	})
	return chengJiaoSpider
}
func (s *ChengJiaoSpider) findAllArea() ([]*model.Area, error) {
	tx := ares.Default().GetOrm("sjz")
	var results []*model.Area

	if err := tx.Model(&model.Area{}).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (s *ChengJiaoSpider) Start() {
	//areas, _ := s.findAllArea()
	//for _, area := range areas {
	//	s.parseOnArea(area)
	//}
	areas, _ := s.findAllArea()
	for _, area := range areas {
		if area.DistrictId == "damacun" {
			fmt.Printf("start parse area: %v\n", area.DistrictId)
			s.parseOnArea(area)
		}
	}
}
func (s *ChengJiaoSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com", ".baidu.com"), //白名单域名
		colly.AllowURLRevisit(),                             //允许对同一 URL 进行多次下载
		colly.Async(true),                                   //设置为异步请求
		colly.Debugger(&debug.LogDebugger{}),                // 开启debug
		colly.MaxDepth(2),                                   //爬取页面深度,最多为两层
		colly.MaxBodySize(1024*1024),                        //响应正文最大字节数
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) "),
		colly.IgnoreRobotsTxt(), //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	//设置请求头
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("My-Header", "dj")
	})

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
func (s *ChengJiaoSpider) parseHouseList(area *model.Area, e *colly.HTMLElement) {
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
