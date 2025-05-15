package zhsc

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis/v8"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
	"time"
)

type PoetrySpider struct {
}

func NewPoetrySpider() *PoetrySpider {
	return &PoetrySpider{}
}

func (s *PoetrySpider) Start() {
	//s.startPoetry(consts.Shi.Name())
	s.startPoetry(consts.Ci.Name())
	//s.startPoetry(consts.Qu.Name())
	//s.startPoetry(consts.Fu.Name())
	//s.startPoetry(consts.Wen.Name())

}
func (s *PoetrySpider) getRandProxy() string {
	authKey := "6E4F6112"
	password := "6028BC3DB221"
	redisKey := "ip_proxy"
	startTime := fmt.Sprintf("%d", time.Now().Unix())
	endTime := fmt.Sprintf("%d", time.Now().Add(1*time.Minute).Unix())
	redisClient := ares.Default().GetRedis("base")
	list := redisClient.ZRangeByScore(context.TODO(), redisKey, &redis.ZRangeBy{
		Min: startTime,
		Max: endTime,
	}).Val()
	if len(list) <= 20 {
		resp, _ := resty.New().SetTimeout(20 * time.Second).SetRetryWaitTime(5 * time.Second).R().Get(fmt.Sprintf("https://share.proxy.qg.net/get?key=%s&pwd=%s", authKey, password))
		ipAddress := gjson.Get(string(resp.Body()), "data.0.server").String()

		if ipAddress != "" {
			_ = redisClient.ZAdd(context.TODO(), redisKey, &redis.Z{
				Score:  cast.ToFloat64(endTime),
				Member: ipAddress,
			}).Val()
			_ = redisClient.ZRemRangeByScore(context.TODO(), redisKey, fmt.Sprintf("%d", 0), startTime).Val()
			list = append(list, ipAddress)
		}
	}
	proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s", authKey, password, lo.Sample(list)))
	return proxyUrl.String()
}
func (s *PoetrySpider) startPoetry(poetryType string) {
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
			time.Sleep(200 * time.Millisecond)
			c.SetProxy(s.getRandProxy())
			c.Visit(fmt.Sprintf("https://zhsc.org%s", hrefStr))
		})
	})
	c.OnHTML(".pagination", func(e *colly.HTMLElement) {
		curPage := cast.ToInt64(strings.TrimSpace(e.DOM.Find(".active").Find("span").Text()))
		var totalPage int64
		e.DOM.Find("li").Find("a").Each(func(i int, item *goquery.Selection) {
			temp := cast.ToInt64(strings.TrimSpace(item.Text()))
			if temp > totalPage {
				totalPage = temp
			}
		})
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			time.Sleep(500 * time.Millisecond)
			c.SetProxy(s.getRandProxy())
			c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-%d.htm", poetryType, curPage+1))
		}
	})
	c.OnResponse(func(r *colly.Response) {
		// 获取当前访问的 URL
		urlStr := r.Request.URL.Path
		if !strings.HasPrefix(urlStr, "/work/work-") {
			return
		}
		poetryId := strings.TrimSuffix(strings.TrimPrefix(urlStr, "/work/work-"), ".htm")
		// 创建文档
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(r.Body)))
		if err != nil {
			fmt.Printf("创建文档失败: %v\n", err)
			return
		}
		s.parsePoetry(poetryId, poetryType, doc.Find(".work"))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
		c.SetProxy(s.getRandProxy())
		c.Visit(r.Request.URL.String())
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.SetProxy(s.getRandProxy())
	c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-2210.htm", poetryType))
}
func (s *PoetrySpider) parsePoetry(poetryId, poetryType string, e *goquery.Selection) {
	poetry := &model.Poetry{}
	poetry.Title = strings.TrimSpace(e.Find(".item-hd").Text())
	poetry.Dynasty = strings.TrimSpace(e.Find(".item-dynasty-author").Find("a").First().Text())
	poetry.AuthorName = strings.TrimSpace(e.Find(".item-dynasty-author").Find("a").Last().Text())
	poetry.Paragraphs = strings.TrimSpace(e.Find(".item-desc.work-content").Text())
	authorId, _ := e.Find(".item-dynasty-author").Find("a").Last().Attr("href")
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
	poetry.PoetryId = poetryId
	interpret := &model.Interpret{}
	interpret.PoetryId = poetry.PoetryId
	interpret.Intro = strings.TrimSpace(e.Find("#intro").Text())
	interpret.Annotation = strings.TrimSpace(e.Find("#annotation").Text())
	interpret.Translation = strings.TrimSpace(e.Find("#translation").Text())

	author := &model.Author{}
	author.AuthorId = poetry.AuthorId
	author.AuthorName = poetry.AuthorName
	author.Dynasty = poetry.Dynasty
	_ = services.GetPoetryService().SavePoetry(poetry, "zhsc_poetry")
	_ = services.GetInterpretService().SaveInterpret(interpret, "zhsc_poetry")
	_ = services.GetAuthorService().SaveAuthor(author, "zhsc_poetry")
}
