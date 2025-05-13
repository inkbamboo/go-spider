package spiders

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders/zhsc"
)

type SpiderInterface interface {
	Start()
}

func NewInstance(platform string) (spider SpiderInterface, err error) {
	if platform == consts.ZHSC.Name() {
		spider = zhsc.NewPoetrySpider()
		return
	}
	err = fmt.Errorf("current platform:%v is not support", platform)
	return
}
