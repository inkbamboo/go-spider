package ke

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/util"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
	"strings"
	"time"
)

type ChengJiaoSpider struct {
	city  string
	alias string
}

func NewChengJiaoSpider(city string) *ChengJiaoSpider {
	return &ChengJiaoSpider{
		city:  city,
		alias: fmt.Sprintf("ke_%s", city),
	}
}

func (s *ChengJiaoSpider) setCookie(c *colly.Collector) {
	//设置请求头
	redisClient := ares.Default().GetRedis("base")
	cookie := redisClient.Get(context.TODO(), "cookie").Val()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", cookie)
	})
}
func (s *ChengJiaoSpider) Start() {
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
func (s *ChengJiaoSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector(
		colly.AllowedDomains(fmt.Sprintf("%s.ke.com", s.city)), //白名单域名
		colly.AllowURLRevisit(),                                //允许对同一 URL 进行多次下载
		colly.Async(true),                                      //设置为异步请求
		colly.MaxDepth(2),                                      //爬取页面深度,最多为两层
		colly.MaxBodySize(1024*1024),                           //响应正文最大字节数
		colly.IgnoreRobotsTxt(),                                //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		RandomDelay: 10 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	//设置请求头
	s.setCookie(c)
	c.OnHTML(".leftContent", func(e *colly.HTMLElement) {
		dealDate := s.parseHouseList(area, e.DOM.Find(".listContent ").Find("li"))
		if !dealDate.IsZero() && dealDate.Before(time.Now().Add(-24*60*time.Hour)) {
			return
		}
		pageData, _ := e.DOM.Find(".page-box div").Attr("page-data")
		totalPage := gjson.Get(pageData, "totalPage").Int()
		curPage := gjson.Get(pageData, "curPage").Int()
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			c.Visit(fmt.Sprintf("https://%s.ke.com/chengjiao/%s/pg%d/", s.city, area.DistrictId, curPage+1))
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://%s.ke.com/chengjiao/%s/", s.city, area.DistrictId))
}
func (s *ChengJiaoSpider) parseHouseList(area *model.Area, e *goquery.Selection) (dealDate time.Time) {
	dealDate = time.Now()
	e.Each(func(i int, el *goquery.Selection) {
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
		}
		titleInfo := strings.TrimSpace(el.Find(".title").Text())
		houseItem.XiaoquName = strings.Split(titleInfo, " ")[0]
		houseItem.HouseArea = util.GetHouseArea(titleInfo)
		houseItem.HouseType = util.GetHouseType(titleInfo)
		houseItem.HouseOrientation = util.GetHouseOrientation(strings.TrimSpace(el.Find(".houseInfo").Text()))
		houseItem.HouseFloor = util.GetHouseFloor(el.Find(".positionInfo").Text())
		chengjiao := &model.ChengJiao{
			HousedelId: housedelId,
			DistrictId: area.DistrictId,
		}
		dealInfo := util.TrimInfoEmpty(el.Find(".dealCycleTxt").Text())
		chengjiao.TotalPrice = util.GetTotalPrice(dealInfo)
		chengjiao.DealCycle = util.GetDealCycle(dealInfo)
		chengjiao.UnitPrice = util.GetUnitPrice(el.Find(".unitPrice").Text())
		chengjiao.DealDate = strings.ReplaceAll(util.TrimInfoEmpty(el.Find(".dealDate").Text()), ".", "-")
		if chengjiao.DealDate != "" {
			curDealDate, _ := time.Parse("2006-01-02", chengjiao.DealDate)
			if !dealDate.IsZero() && curDealDate.Before(dealDate) {
				dealDate = curDealDate
			}
		}
		chengjiao.DealPrice = util.GetTotalPrice(util.TrimInfoEmpty(el.Find(".totalPrice").Text()))
		if err := services.GetHouseService().SaveHouse(houseItem, s.alias); err != nil {
			return
		}
		if err := services.GetChengJiaoService().SaveChengJiao(chengjiao, s.alias); err != nil {
			return
		}
	})
	return
}
