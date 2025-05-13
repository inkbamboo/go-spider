package zhsc

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/util"
	"github.com/spf13/cast"
	"strings"
	"time"
)

type PoetrySpider struct {
}

func NewPoetrySpider() *PoetrySpider {
	return &PoetrySpider{}
}

func (s *PoetrySpider) Start() {
	//https://zhsc.org/shi/page-2.htm
	s.parsePoetry(consts.Shi.Name())
	//s.parsePoetry(consts.Ci.Name())
	//s.parsePoetry(consts.Qu.Name())
	//s.parsePoetry(consts.Fu.Name())
	//s.parsePoetry(consts.Wen.Name())

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
	c.OnHTML(".item-list", func(e *colly.HTMLElement) {
		e.DOM.Find(".item-btn").Each(func(i int, selection *goquery.Selection) {
			hrefStr, _ := selection.Attr("href")
			if hrefStr == "" {
				return
			}
			c.Visit(fmt.Sprintf("https://zhsc.org%s", hrefStr))
		})
	})
	c.OnHTML(".pagination", func(e *colly.HTMLElement) {
		curPage := cast.ToInt64(strings.TrimSpace(e.DOM.Find(".active").Find("span").Text()))
		totalPage := cast.ToInt64(strings.TrimSpace(e.DOM.Find("li").Find("a").Last().Text()))
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			time.Sleep(200 * time.Millisecond)
			c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-%d.htm", poetryType, curPage+1))
		}
	})

	c.OnHTML(".work", func(e *colly.HTMLElement) {
		poetry := &model.Poetry{}
		poetry.Title = strings.TrimSpace(e.DOM.Find(".item-hd").Text())
		poetry.Dynasty = strings.TrimSpace(e.DOM.Find(".item-dynasty-author").Find("a").First().Text())
		poetry.Author = strings.TrimSpace(e.DOM.Find(".item-dynasty-author").Find("a").Last().Text())
		poetry.Paragraphs = strings.TrimSpace(e.DOM.Find(".item-desc.work-content").Text())
		authorId, _ := e.DOM.Find(".item-dynasty-author").Find("a").Last().Attr("href")
		poetry.AuthorId = strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(authorId), "/author/author-"), ".htm")
		switch poetryType {
		case consts.Shi.Name():
			poetry.PoetryType = consts.Shi.Description()
		case consts.Ci.Name():
			poetry.PoetryType = consts.Ci.Description()
		case consts.Qu.Name():
			poetry.PoetryType = consts.Qu.Description()
		case consts.Wen.Name():
			poetry.PoetryType = consts.Wen.Description()
		case consts.Fu.Name():
			poetry.PoetryType = consts.Fu.Description()
		}
		poetry.PoetryId = util.GetMd5(poetry.Paragraphs)
		interpret := &model.Interpret{}
		interpret.PoetryId = poetry.PoetryId
		interpret.Intro = strings.TrimSpace(e.DOM.Find("#intro").Text())
		interpret.Annotation = strings.TrimSpace(e.DOM.Find("#annotation").Text())
		interpret.Translation = strings.TrimSpace(e.DOM.Find("#translation").Text())

		author := &model.Author{}
		author.AuthorId = poetry.AuthorId
		author.Name = poetry.Author
		author.Dynasty = poetry.Dynasty
		_ = services.GetPoetryService().SavePoetry(poetry, "zhsc_poetry")
		_ = services.GetInterpretService().SaveInterpret(interpret, "zhsc_poetry")
		_ = services.GetAuthorService().SaveAuthor(author, "zhsc_poetry")
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-150.htm", poetryType))
}
