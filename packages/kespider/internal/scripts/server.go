package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spiders"
	"time"
)

func RunAreaSpider() {
	fmt.Println("RunAreaSpider")
	spiders.GetAreaSpider().Start()
	time.Sleep(100 * time.Second)
}
func RunErShouSpider() {
	spiders.GetErShouSpider().Start()
	//crontab := cron.New(cron.WithSeconds())
	//if _, err := crontab.AddFunc(fmt.Sprintf("0 3 1-31/5 * * *"), spider.GetErShouSpider().Start); err != nil {
	//	fmt.Println("GetAreaSpider err: ", err)
	//}
	select {}
}
