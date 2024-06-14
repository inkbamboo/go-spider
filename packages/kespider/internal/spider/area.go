package spider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	//c.Limit(&colly.LimitRule{
	//	DomainGlob:  "sjz.ke.com/*",
	//	Delay:       1 * time.Second,
	//	RandomDelay: 1 * time.Second,
	//	Parallelism: 0,
	//})
	// Find and visit all links
	//c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	//	e.Request.Visit(e.Attr("href"))
	//})

	c.OnResponse(func(r *colly.Response) {
		doc, err := goquery.NewDocument(string(r.Body))
		fmt.Println("Visited", err)
		fmt.Println("Visited", doc.Text())

	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://sjz.ke.com")
}
