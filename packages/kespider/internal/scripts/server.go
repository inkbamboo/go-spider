package scripts

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/model"
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spiders"
	"github.com/samber/lo"
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
	tx := ares.Default().GetOrm("sjz")
	versions := []string{"2024-06-25", "2024-06-20"}
	var hosePrices []*model.HousePrice
	_ = tx.Model(&model.HousePrice{}).Where("version in(?)", versions).Find(&hosePrices).Error
	houseInfos := lo.GroupBy(hosePrices, func(item *model.HousePrice) string {
		return item.HousedelId
	})
	var sellOutHouse, newHouse, changeHouse string
	for _, houseInfo := range houseInfos {
		var oldPrice, newPrice float64
		for _, item := range houseInfo {
			if item.Version == "2024-06-20" {
				oldPrice = item.TotalPrice
			} else if item.Version == "2024-06-25" {
				newPrice = item.TotalPrice
			}
		}
		if newPrice == 0 {
			sellOutHouse += fmt.Sprintf("%s\n", houseInfo[0].HousedelId)
			continue
		}
		if oldPrice == 0 {
			newHouse += fmt.Sprintf("%s\n", houseInfo[0].HousedelId)
			continue
		}
		if newPrice < oldPrice {
			changeHouse += fmt.Sprintf("%s,%v,%v\n", houseInfo[0].HousedelId, oldPrice, newPrice)
			continue
		}
	}
	os.WriteFile("sell_out_house.txt", []byte(sellOutHouse), 0644)
	os.WriteFile("new_house.txt", []byte(newHouse), 0644)
	os.WriteFile("change_house.txt", []byte(changeHouse), 0644)
	select {}
}
