package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders"
	"time"
)

func RunPoetrySpider(platform string) {
	fmt.Println("RunPoetrySpider")
	if sp, err := spiders.NewInstance(platform); err != nil {
		fmt.Println("RunPoetrySpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
	time.Sleep(100 * time.Second)
}
func RunTest() {
	fmt.Printf("*********%+v\n", consts.PoetryTypeNames())
	fmt.Printf("*********%+v\n", consts.PoetryTypeNames())

}
