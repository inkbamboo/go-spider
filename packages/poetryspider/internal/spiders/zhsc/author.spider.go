package zhsc

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"time"
)

type AuthorSpider struct {
}

func NewAuthorSpider() *AuthorSpider {
	return &AuthorSpider{}
}

func (s *AuthorSpider) Start() {
	s.parseAuthor(consts.Shi.Name())
}
func (s *AuthorSpider) parseAuthor(poetryType string) {
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
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	authors, _ := services.GetAuthorService().GetAllAuthor("zhsc_poetry")
	for _, item := range authors {
		time.Sleep(50 * time.Millisecond)
		c.Visit(fmt.Sprintf("https://zhsc.org/author/author-%s.htm", item.AuthorId))

	}
}
