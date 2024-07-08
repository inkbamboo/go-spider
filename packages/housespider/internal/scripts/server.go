package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/housespider/consts"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/services"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/spiders"
	"os"
	"time"
)

func RunAreaSpider(platform, city string) {
	fmt.Println("RunAreaSpider")
	if sp, err := spiders.NewInstance(platform, city, consts.Area.Name()); err != nil {
		fmt.Println("GetAreaSpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
	time.Sleep(100 * time.Second)
}
func RunErShouSpider(platform, city string) {
	if sp, err := spiders.NewInstance(platform, city, consts.ErShou.Name()); err != nil {
		fmt.Println("GetErshouSpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
}
func RunChengJiaoSpider(platform, city string) {
	if sp, err := spiders.NewInstance(platform, city, consts.ChengJiao.Name()); err != nil {
		fmt.Println("GetChengJiaoSpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
}

func RunTest() {
	fmt.Println("RunTest")
	versions := []string{"2024-06-26", "2024-06-30", "2024-07-03", "2024-07-07"}
	changeHouse := services.GetHousePriceService().GetChangeHouse(versions)
	os.WriteFile("./temp/change_house.csv", []byte(changeHouse), 0644)
	time.Sleep(10 * time.Second)
	fmt.Println("RunTest End")
}
