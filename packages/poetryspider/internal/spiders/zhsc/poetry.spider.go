package zhsc

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/tidwall/gjson"
	"time"
)

type PoetrySpider struct {
}

func NewPoetrySpider() *PoetrySpider {
	return &PoetrySpider{}
}

func (s *PoetrySpider) Start() {
	type SpiderTypeEnum struct {
		Shi string `enum:"shi,诗"`
		Ci  string `enum:"ci,词"`
		Wen string `enum:"wen,文"`
		Qu  string `enum:"qu,曲"`
		Fu  string `enum:"fu,赋"`
	}
	//https://zhsc.org/shi/page-2.htm
	s.parsePoetry("shi")
}
func (s *PoetrySpider) parsePoetry(poetryType string) {
	c := colly.NewCollector(
		colly.AllowedDomains("zhsc.org"), //白名单域名
		colly.AllowURLRevisit(),          //允许对同一 URL 进行多次下载
		colly.Async(true),                //设置为异步请求
		colly.MaxDepth(2),                //爬取页面深度,最多为两层
		colly.MaxBodySize(2*1024*1024),   //响应正文最大字节数
		colly.IgnoreRobotsTxt(),          //忽略目标机器中的`robots.txt`声明
	)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		RandomDelay: 10 * time.Second})
	//随机设置User-Agent
	extensions.RandomUserAgent(c)
	c.OnHTML("bdy", func(e *colly.HTMLElement) {
		//s.parseHouseList(area, e.DOM.Find(".sellListContent").Find("li"))
		pageData, _ := e.DOM.Find(".pagination").Attr("list")
		totalPage := gjson.Get(pageData, "totalPage").Int()
		curPage := gjson.Get(pageData, "curPage").Int()
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			//c.Visit(fmt.Sprintf("https://%s.ke.com/ershoufang/%s/pg%d/", s.city, area.DistrictId, curPage+1))
		}
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-1.htm", poetryType))
}
