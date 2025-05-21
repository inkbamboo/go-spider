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
	cache    *cache.Cache
	checkUrl string
	endPage  int64
}

func NewPoetrySpider() *PoetrySpider {
	return &PoetrySpider{
		cache: cache.New(time.Minute*1, time.Second*10),
	}
}

func (s *PoetrySpider) Start() {
	//s.startPoetry(consts.Shi.Name())
	s.startPoetry(consts.Ci.Name())
	//s.startPoetry(consts.Qu.Name())
	//s.startPoetry(consts.Fu.Name())
	//s.startPoetry(consts.Wen.Name())

}
func (s *PoetrySpider) getLocalCacheProxy() (proxyList []*ProxyInfo) {
	if cacheList, ok := s.cache.Get(consts.GetProxyListKey(ares.GetEnv())); ok {
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
	if cacheList := redisClient.Get(context.TODO(), consts.GetProxyListKey(ares.GetEnv())).Val(); len(cacheList) > 0 {
		_ = sonic.UnmarshalString(cacheList, &proxyList)

	}
	proxyList = lo.FilterMap(proxyList, func(item *ProxyInfo, index int) (*ProxyInfo, bool) {
		expireTime, _ := time.ParseInLocation("2006-01-02 15:04:05", item.ExpireTime, time.Local)
		return item, expireTime.After(time.Now())
	})
	return
}
func (s *PoetrySpider) setCacheProxy(proxyList []*ProxyInfo) {
	if len(proxyList) == 0 {
		s.cache.Set(consts.GetHasProxyKey(ares.GetEnv()), true, cache.DefaultExpiration)
		return
	}
	redisClient := ares.Default().GetRedis("base")
	s.cache.Set(consts.GetProxyListKey(ares.GetEnv()), proxyList, cache.DefaultExpiration)
	str, _ := sonic.MarshalString(proxyList)
	_ = redisClient.Set(context.TODO(), consts.GetProxyListKey(ares.GetEnv()), str, 0).Val()
	return
}
func (s *PoetrySpider) checkProxy(proxyList []*ProxyInfo) []*ProxyInfo {
	var resultList []*ProxyInfo
	for _, proxyInfo := range proxyList {
		proxyUrl, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", proxyInfo.Username, proxyInfo.Password, proxyInfo.Ip, proxyInfo.Port))
		// 发送请求
		resp, err := resty.New().SetProxy(proxyUrl.String()).R().Get(s.checkUrl)
		if err != nil {
			fmt.Printf("Error: %v  proxyUrl %v\n", err, proxyUrl.String())
			continue
		}
		respBody := string(resp.Body())
		fmt.Println("Response Body:", respBody)
		resultList = append(resultList, proxyInfo)
	}
	return resultList
}
func (s *PoetrySpider) getProxyList() (proxyList []*ProxyInfo) {
	if hasProxy, ok := s.cache.Get(consts.GetHasProxyKey(ares.GetEnv())); ok && !hasProxy.(bool) {
		return
	}
	if proxyList = s.getLocalCacheProxy(); len(proxyList) > 0 {
		return
	}
	if proxyList = s.getRedisProxy(); len(proxyList) > 0 {
		return
	}
	redisClient := ares.Default().GetRedis("base")
	if locked := redisClient.SetNX(context.TODO(), consts.GetProxyListLockKey(ares.GetEnv()), 1, 300*time.Second).Val(); !locked {
		return
	}
	resp, _ := resty.New().SetTimeout(20 * time.Second).SetRetryWaitTime(5 * time.Second).R().Get(ares.GetConfig().GetString("proxy"))
	respBody := string(resp.Body())
	_ = sonic.UnmarshalString(gjson.Get(respBody, "data").String(), &proxyList)
	fmt.Printf("getProxyList%v\n", respBody)
	if len(proxyList) > 0 {
		if proxyList = s.checkProxy(proxyList); len(proxyList) == 0 {
			_ = redisClient.Del(context.TODO(), consts.GetProxyListLockKey(ares.GetEnv())).Val()
			return s.getProxyList()
		}
	}
	s.setCacheProxy(proxyList)
	_ = redisClient.Del(context.TODO(), consts.GetProxyListLockKey(ares.GetEnv())).Val()
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

func (s *PoetrySpider) getPoetryPage(poetryType, pageType string) (page int64) {
	redisClient := ares.Default().GetRedis("base")
	page, _ = redisClient.Get(context.TODO(), consts.GetPoetryPageKey(ares.GetEnv(), poetryType, pageType)).Int64()
	if page == 0 {
		page = 1
	}
	return
}
func (s *PoetrySpider) setPoetryPage(poetryType, pageType string, page int64) {
	redisClient := ares.Default().GetRedis("base")
	redisClient.Set(context.TODO(), consts.GetPoetryPageKey(ares.GetEnv(), poetryType, pageType), page, 0).Val()
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
	c.SetRequestTimeout(time.Duration(ares.GetConfig().GetInt64("timeout")) * time.Second)
	_ = c.Limit(&colly.LimitRule{
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
			poetryId := strings.TrimSuffix(strings.TrimPrefix(hrefStr, "/work/work-"), ".htm")
			if !services.GetPoetryService().PoetryExists(poetryId, "zhsc_poetry") {
				_ = c.SetProxy(s.getRandProxy())
				_ = c.Visit(fmt.Sprintf("https://zhsc.org%s", hrefStr))
			}

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
		if s.endPage > 0 && s.endPage < totalPage {
			totalPage = s.endPage
		}
		if curPage < totalPage {
			c.UserAgent = browser.Random()
			time.Sleep(100 * time.Millisecond)
			urlStr := fmt.Sprintf("https://zhsc.org/%s/page-%d.htm", poetryType, curPage+1)
			s.checkUrl = urlStr
			_ = c.SetProxy(s.getRandProxy())
			_ = c.Visit(urlStr)
		}
	})
	c.OnResponse(func(r *colly.Response) {
		// 获取当前访问的 URL
		urlStr := r.Request.URL.Path
		fmt.Println("Visiting", urlStr)
		if strings.HasPrefix(urlStr, fmt.Sprintf("/%s/page-", poetryType)) {
			currentPage := strings.TrimSuffix(strings.TrimPrefix(urlStr, fmt.Sprintf("/%s/page-", poetryType)), ".htm")
			s.setPoetryPage(poetryType, "start", cast.ToInt64(currentPage))
			return
		}
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
		_ = c.SetProxy(s.getRandProxy())
		_ = c.Visit(r.Request.URL.String())
	})
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting", r.URL)
	})
	_ = c.SetProxy(s.getRandProxy())
	//还差 3160-4500
	//s.endPage = 4500
	//urlStr := fmt.Sprintf("https://zhsc.org/%s/page-3160.htm", poetryType)
	//从 1 开始
	//urlStr:=fmt.Sprintf("https://zhsc.org/%s/page-3956.htm", poetryType)

	// 从 60000 开始
	//s.endPage = 72800
	//urlStr := fmt.Sprintf("https://zhsc.org/%s/page-67003.htm", poetryType)
	// 从 72800 开始
	startPage := s.getPoetryPage(poetryType, "start")
	s.endPage = s.getPoetryPage(poetryType, "end")
	urlStr := fmt.Sprintf("https://zhsc.org/%s/page-%d.htm", poetryType, startPage)
	s.checkUrl = urlStr
	_ = c.Visit(urlStr)
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
