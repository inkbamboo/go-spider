package scripts

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders"
)

func RunPoetrySpider(platform, spider string) {
	fmt.Println("RunPoetrySpider")
	if sp, err := spiders.NewInstance(platform, spider); err != nil {
		fmt.Println("RunPoetrySpider err: ", err)
	} else {
		sp.Start()
	}
	select {}
}
func RunTest() {
	fmt.Printf("*********%+v\n", consts.PoetryTypeNames())
	fmt.Printf("*********%+v\n", consts.PoetryTypeNames())

}
