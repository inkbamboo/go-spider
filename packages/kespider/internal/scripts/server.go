package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spider"
	"time"
)

func RunAreaSpider() {
	fmt.Println("RunAreaSpider")
	spider.GetAreaSpider().Start()
	time.Sleep(100 * time.Second)
}
func RunErShouSpider() {
	spider.GetErShouSpider().Start()
	//crontab := cron.New(cron.WithSeconds())
	//if _, err := crontab.AddFunc(fmt.Sprintf("0 3 1-31/5 * * *"), spider.GetErShouSpider().Start); err != nil {
	//	fmt.Println("GetAreaSpider err: ", err)
	//}
	select {}
}
func RunTest() {

}
