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
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", "lianjia_uuid=ca7c2641-e6ed-46eb-9c0f-fd592d9af6e5; digv_extends=%7B%22utmTrackId%22%3A%22%22%7D; crosSdkDT2019DeviceId=rcz5ws-mye2yu-rdk001maa3w0kf9-0e7zpslu8; _ga=GA1.2.1425795671.1713798460; ke_uuid=7de6895ca08a7f4e9b0c379251e2186b; lfrc_=e61f584b-111f-415b-9004-75e33b7e1677; __xsptplus788=788.4.1713947645.1713947645.1%234%7C%7C%7C%7C%7C%23%23%23; Qs_lvt_200116=1718959950; Qs_pv_200116=13851011825112360; HMACCOUNT=7F7556E99A058C99; lianjia_ssid=2e80527d-9662-49f2-8671-95be1dee624d; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2218f0658f9c7d56-03a07871dec8a7-1b525637-1484784-18f0658f9c81e11%22%2C%22%24device_id%22%3A%2218f0658f9c7d56-03a07871dec8a7-1b525637-1484784-18f0658f9c81e11%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E4%BB%98%E8%B4%B9%E5%B9%BF%E5%91%8A%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_utm_source%22%3A%22baidu%22%2C%22%24latest_utm_medium%22%3A%22pinzhuan%22%2C%22%24latest_utm_campaign%22%3A%22wybeijing%22%2C%22%24latest_utm_content%22%3A%22biaotimiaoshu%22%2C%22%24latest_utm_term%22%3A%22biaoti%22%7D%7D; select_city=130100; Hm_lvt_9152f8221cb6243a53c83b956842be8a=1722392537; Hm_lpvt_9152f8221cb6243a53c83b956842be8a=1722392900; srcid=eyJ0Ijoie1wiZGF0YVwiOlwiMzQ4Nzk0M2I2OGVmNTkyNGJjYmMyZjljYWNkMTRlYzcxNGM5YTA1OWRkMzM0YTA2NDI5NWI4NmIzYWRhOTY0ZGIzZTQxZTgyZjg0ZGJjMGFlNzY0MTc0NTlkMWQ3MjhlYmU4ZDllYmI2ZmU4ZDVhYjg1Nzc3YWIzNDRkZWRlZmZlOTM1M2NjYzVlMjdkMzM4NjRkZGI1ZTUwYTgwN2FiZTgzYmQxMjExM2MyNzJlZDBhYmQ3Yjk2Y2ZmNDViZGNlOTkwNTkxYmUxMTA1YWQ2ODZiMzAxYWY2MDk2MzgyZDU2MzA0OTdlMjQzOTYxNzc5OTA4NmE5OTcxOTM2ZWY2ZlwiLFwia2V5X2lkXCI6XCIxXCIsXCJzaWduXCI6XCJjYmFhZWVmNlwifSIsInIiOiJodHRwczovL3Nqei5rZS5jb20vIiwib3MiOiJ3ZWIiLCJ2IjoiMC4xIn0=; login_ucid=2000000109247156; lianjia_token=2.001153fbc2724a6c9a00fed2f3857b2c3a; lianjia_token_secure=2.001153fbc2724a6c9a00fed2f3857b2c3a; security_ticket=J2kYMdYuPv9uqbr97ixI6kcVQBnyMY5a246iPxcy8xbi6nKp9Tw4532k16L/t96RXRp+1tg1kjx5Izf+iARW9Kxj38pNoDXmaFi83vRE4789Pxj2UvvLey7/5cLSyCWhz+jOZUHrsMF2Nv4P1l13rBdnrUECiyyxrrln3Y6T1H4=; ftkrc_=b78807e4-23bd-495a-928f-0cadff56bc27")
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
