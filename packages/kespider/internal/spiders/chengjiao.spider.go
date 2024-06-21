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

func (s *ChengJiaoSpider) setCookie(c *colly.Collector) {
	//设置请求头
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", "lianjia_uuid=ca7c2641-e6ed-46eb-9c0f-fd592d9af6e5; digv_extends=%7B%22utmTrackId%22%3A%22%22%7D; crosSdkDT2019DeviceId=rcz5ws-mye2yu-rdk001maa3w0kf9-0e7zpslu8; _ga=GA1.2.1425795671.1713798460; ke_uuid=7de6895ca08a7f4e9b0c379251e2186b; ftkrc_=4bb9ecef-5a80-4e2c-adf6-ab0468a8b9a0; lfrc_=e61f584b-111f-415b-9004-75e33b7e1677; __xsptplus788=788.4.1713947645.1713947645.1%234%7C%7C%7C%7C%7C%23%23%23; Hm_lvt_9152f8221cb6243a53c83b956842be8a=1718709243; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2218f0658f9c7d56-03a07871dec8a7-1b525637-1484784-18f0658f9c81e11%22%2C%22%24device_id%22%3A%2218f0658f9c7d56-03a07871dec8a7-1b525637-1484784-18f0658f9c81e11%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E4%BB%98%E8%B4%B9%E5%B9%BF%E5%91%8A%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22https%3A%2F%2Fwww.baidu.com%2Fother.php%22%2C%22%24latest_referrer_host%22%3A%22www.baidu.com%22%2C%22%24latest_search_keyword%22%3A%22%E8%B4%9D%E5%A3%B3%E6%89%BE%E6%88%BF%22%2C%22%24latest_utm_source%22%3A%22baidu%22%2C%22%24latest_utm_medium%22%3A%22pinzhuan%22%2C%22%24latest_utm_campaign%22%3A%22wybeijing%22%2C%22%24latest_utm_content%22%3A%22biaotimiaoshu%22%2C%22%24latest_utm_term%22%3A%22biaoti%22%7D%7D; lianjia_ssid=1748ec9d-8cbb-46e5-a333-f093c8ff69a1; login_ucid=2000000109247156; lianjia_token=2.0012310297712895cf039c2ba6757035ba; lianjia_token_secure=2.0012310297712895cf039c2ba6757035ba; security_ticket=bQlnrwkSU7dXYfQbMLSfgaMQFzH+hHbjfQswMuCuZdoghSi9MedMLNKRZfOmMedG40ZaK3RMBj5iIyIj3txfd9k8J6Qm1jcSMWY7YUd6NSQHe0gqdoM9O9/hzPRmu4Ia7SzM/LPGJ0ovgj+hAUadTeYjmv8aAS6z63z5j+mgZ7g=; select_city=130100; Hm_lpvt_9152f8221cb6243a53c83b956842be8a=1718947216; srcid=eyJ0Ijoie1wiZGF0YVwiOlwiZWUyODIzNzE1NzMyNWQxNWRmNzllN2MzMjdkYWM4ZjRkNzBlMzIzN2E3N2U2ZDgwODIzM2RlMDU1OGU0NzczYjg3NDQ4YTIwMWE1MTA2MTI0OWQ1MjdmMmU0YjQ0NGZiZDVmZjQ0MTg5MjdiYTMwZmI5Mzc3ODYyYWJiYjJiMDI0MDk0MGU0ZWQyZmNhMGI5ZDY4N2I5OWQ0ZjJkNjYzNGE2NmE5M2EzM2VkODg0MDViOGM0NjI4ZGNkMmY1YTA1MWY1YzFlZjlmNmFiNjg0NjRkNGZkYWYxMTA5YmMwYjlmZjhjMDM2MTM1MWE0OTdmM2FjZjc4NmU4YjgwNWJjOVwiLFwia2V5X2lkXCI6XCIxXCIsXCJzaWduXCI6XCIzMGYyZDI5Y1wifSIsInIiOiJodHRwczovL3Nqei5rZS5jb20vY2hlbmdqaWFvL2NoYW5nYW4vcGcyLyIsIm9zIjoid2ViIiwidiI6IjAuMSJ9")
	})
}
func (s *ChengJiaoSpider) Start() {
	areas, _ := services.GetAreaService().FindAllArea()
	for _, area := range areas {
		s.parseOnArea(area)
	}
	//areas, _ := services.GetAreaService().FindAllArea()
	//for _, area := range areas {
	//	if area.DistrictId == "damacun" {
	//		fmt.Printf("start parse area: %v\n", area.DistrictId)
	//		s.parseOnArea(area)
	//	}
	//}
}
func (s *ChengJiaoSpider) parseOnArea(area *model.Area) {
	c := colly.NewCollector(
		colly.AllowedDomains("sjz.ke.com"), //白名单域名
		colly.AllowURLRevisit(),            //允许对同一 URL 进行多次下载
		colly.Async(true),                  //设置为异步请求
		colly.MaxDepth(2),                  //爬取页面深度,最多为两层
		colly.MaxBodySize(1024*1024),       //响应正文最大字节数
		colly.IgnoreRobotsTxt(),            //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	//设置请求头
	s.setCookie(c)
	c.OnHTML(".listContent ", func(e *colly.HTMLElement) {
		s.parseHouseList(area, e)
	})
	c.OnHTML(".page-box div", func(e *colly.HTMLElement) {
		totalPage := gjson.Get(e.Attr("page-data"), "totalPage").Int()
		curPage := gjson.Get(e.Attr("page-data"), "curPage").Int()
		fmt.Printf("totalPage: %d curPage: %d\n", totalPage, curPage)
		//if curPage < totalPage {
		//	c.UserAgent = ""
		//	c.Visit(fmt.Sprintf("https://sjz.ke.com/chengjiao/%s/pg%d/", area.DistrictId, curPage+1))
		//}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://sjz.ke.com/chengjiao/%s/", area.DistrictId))
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
		}
		titleInfo := strings.TrimSpace(el.DOM.Find(".title").Text())
		houseItem.XiaoquName = strings.Split(titleInfo, " ")[0]
		houseItem.HouseArea = util.GetHouseArea(titleInfo)
		houseItem.HouseType = util.GetHouseType(titleInfo)
		houseItem.HouseOrientation = util.GetHouseOrientation(strings.TrimSpace(el.DOM.Find(".houseInfo").Text()))
		houseItem.HouseFloor = util.GetHouseFloor(el.DOM.Find(".positionInfo").Text())
		chengjiao := &model.ChengJiao{
			HousedelId: housedelId,
			DistrictId: area.DistrictId,
		}
		dealInfo := strings.TrimSpace(el.DOM.Find(".dealCycleTxt").Text())

		chengjiao.TotalPrice = util.GetTotalPrice(dealInfo)
		chengjiao.DealCycle = util.GetDealCycle(dealInfo)
		chengjiao.UnitPrice = util.GetUnitPrice(el.DOM.Find(".unitPrice").Text())
		chengjiao.DealDate = strings.ReplaceAll(util.TrimInfoEmpty(el.DOM.Find(".dealDate").Text()), ".", "-")
		chengjiao.DealPrice = util.GetTotalPrice(util.TrimInfoEmpty(el.DOM.Find(".totalPrice").Text()))
		if err := services.GetHouseService().SaveHouse(houseItem); err != nil {
			return
		}
		if err := services.GetChengJiaoService().SaveChengJiao(chengjiao); err != nil {
			return
		}
	})
}
