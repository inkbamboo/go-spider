package zhsc

import (
	"context"
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/bytedance/sonic"
	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/model"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/services"
	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"net/url"
	"strings"
	"time"
)

type ProxyInfo struct {
	Ip         string `json:"ip"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Ttl        int    `json:"ttl"`
	ExpireTime string `json:"expireTime"`
}
type PoetrySpider struct {
	cache *cache.Cache
}

func NewPoetrySpider() *PoetrySpider {
	return &PoetrySpider{
		cache: cache.New(time.Minute*5, time.Minute*1),
	}
}

func (s *PoetrySpider) Start() {
	s.startPoetry(consts.Shi.Name())
	//s.startPoetry(consts.Ci.Name())
	//s.startPoetry(consts.Qu.Name())
	//s.startPoetry(consts.Fu.Name())
	//s.startPoetry(consts.Wen.Name())

}
func (s *PoetrySpider) getLocalCacheProxy() (proxyList []*ProxyInfo) {
	redisKey := "proxy_list"
	if cacheList, ok := s.cache.Get(redisKey); ok {
		proxyList = cacheList.([]*ProxyInfo)
		proxyList = lo.FilterMap(proxyList, func(item *ProxyInfo, index int) (*ProxyInfo, bool) {
			expireTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.ExpireTime, time.Local)
			return item, expireTime.After(time.Now())
		})
	}
	return
}
func (s *PoetrySpider) getRedisProxy() (proxyList []*ProxyInfo) {
	redisClient := ares.Default().GetRedis("base")
	redisKey := "proxy_list"
	if cacheList := redisClient.Get(context.TODO(), redisKey).String(); len(cacheList) > 0 {
		_ = sonic.UnmarshalString(cacheList, &proxyList)
	}
	proxyList = lo.FilterMap(proxyList, func(item *ProxyInfo, index int) (*ProxyInfo, bool) {
		return item, item.ExpireTime.After(time.Now())
	})
	return
}
func (s *PoetrySpider) setCacheProxy(proxyList []*ProxyInfo) {
	if len(proxyList) == 0 {
		s.cache.Set("hasProxy", true, cache.DefaultExpiration)
		return
	}
	redisClient := ares.Default().GetRedis("base")
	redisKey := "proxy_list"
	s.cache.Set(redisKey, proxyList, cache.DefaultExpiration)
	str, _ := sonic.MarshalString(proxyList)
	_ = redisClient.Set(context.TODO(), redisKey, str, 0).Val()
	return
}
func (s *PoetrySpider) getProxyList() (proxyList []*ProxyInfo) {
	if hasProxy, ok := s.cache.Get("hasProxy"); ok && !hasProxy.(bool) {
		return
	}
	if proxyList = s.getLocalCacheProxy(); len(proxyList) > 0 {
		return
	}
	if proxyList = s.getRedisProxy(); len(proxyList) > 0 {
		return
	}
	resp, _ := resty.New().SetTimeout(20 * time.Second).SetRetryWaitTime(5 * time.Second).R().Get(ares.GetConfig().GetString("proxy"))
	respBody := string(resp.Body())
	_ = sonic.UnmarshalString(gjson.Get(respBody, "data").String(), &proxyList)
	fmt.Printf("*******respBody start %v %v\n", gjson.Get(respBody, "data").String(), proxyList)

	s.setCacheProxy(proxyList)
	return
}
func (s *PoetrySpider) getRandProxy() (proxyUrl string) {
	proxyInfo := lo.Sample(s.getProxyList())
	if proxyInfo == nil {
		return
	}
	proxy, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", proxyInfo.Username, proxyInfo.Password, proxyInfo.Ip, proxyInfo.Port))
	proxyUrl = proxy.String()
	return
}
func (s *PoetrySpider) startPoetry(poetryType string) {
	c := colly.NewCollector(
		colly.AllowedDomains("zhsc.org"), //白名单域名
		colly.AllowURLRevisit(),          //允许对同一 URL 进行多次下载
		colly.Async(true),                //设置为异步请求
		colly.MaxDepth(2),                //爬取页面深度,最多为两层
		colly.MaxBodySize(100*1024*1024), //响应正文最大字节数
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
			time.Sleep(50 * time.Millisecond)
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
			time.Sleep(100 * time.Millisecond)
			c.SetProxy(s.getRandProxy())
			c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-%d.htm", poetryType, curPage+1))
		}
	})
	c.OnResponse(func(r *colly.Response) {
		// 获取当前访问的 URL
		urlStr := r.Request.URL.Path
		fmt.Println("Visiting", urlStr)
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
		//fmt.Println("Request URL:", r.Request.URL, "\nError:", err)
		c.SetProxy(s.getRandProxy())
		c.Visit(r.Request.URL.String())
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})
	c.SetProxy(s.getRandProxy())
	//还差 3160-4500
	//c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-3160.htm", poetryType))

	// 诗 1-3956 已有  72800 到当前程序爬的位置
	//c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-3956.htm", poetryType))
	//c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-72800.htm", poetryType))
	c.Visit(fmt.Sprintf("https://zhsc.org/%s/page-74779.htm", poetryType))
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
