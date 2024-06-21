package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spiders"
	"time"
)

func RunAreaSpider(city string) {
	fmt.Println("RunAreaSpider")
	spiders.NewAreaSpider(city).Start()
	time.Sleep(100 * time.Second)
}
func RunErShouSpider(city string) {
	spiders.NewErShouSpider(city).Start()
	//crontab := cron.New(cron.WithSeconds())
	//if _, err := crontab.AddFunc(fmt.Sprintf("0 3 1-31/5 * * *"), spider.GetErShouSpider().Start); err != nil {
	//	fmt.Println("GetAreaSpider err: ", err)
	//}
	select {}
}
func RunChengJiaoSpider(city string) {
	spiders.NewChengJiaoSpider(city).Start()
	//crontab := cron.New(cron.WithSeconds())
	//if _, err := crontab.AddFunc(fmt.Sprintf("0 3 1-31/5 * * *"), spider.GetErShouSpider().Start); err != nil {
	//	fmt.Println("GetAreaSpider err: ", err)
	//}
	select {}
}
