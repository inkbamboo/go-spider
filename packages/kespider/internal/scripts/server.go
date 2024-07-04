package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spiders"
	"os"
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

func RunTest() {
	fmt.Println("RunTest")
	versions := []string{"2024-06-26", "2024-06-30", "2024-07-03"}
	changeHouse := services.GetHousePriceService().GetChangeHouse(versions)
	os.WriteFile("./temp/change_house.csv", []byte(changeHouse), 0644)
	time.Sleep(10 * time.Second)
	fmt.Println("RunTest End")
}
