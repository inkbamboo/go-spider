package spiders

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/poetryspider/consts"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/spiders/zhsc"
)

type SpiderInterface interface {
	Start()
}

func NewInstance(platform, spiderType string) (spider SpiderInterface, err error) {
	if platform == consts.ZHSC.Name() {
		switch spiderType {
		case consts.Author.Name():

		case consts.Poetry.Name():
			spider = zhsc.NewPoetrySpider()
		}
		return
	}
	err = fmt.Errorf("current platform:%v is not support", platform)
	return
}
