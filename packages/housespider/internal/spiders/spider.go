package spiders

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/housespider/consts"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/spiders/anjuke"
	"github.com/inkbamboo/go-spider/packages/housespider/internal/spiders/ke"
)

type SpiderInterface interface {
	Start()
}

func NewInstance(platform, city, spiderType string) (spider SpiderInterface, err error) {
	if platform == consts.Ke.Name() {
		switch spiderType {
		case consts.Area.Name():
			spider = ke.NewAreaSpider(city)
		case consts.ErShou.Name():
			spider = ke.NewErShouSpider(city)
		case consts.ChengJiao.Name():
			spider = ke.NewChengJiaoSpider(city)
		default:
			err = fmt.Errorf("current platform:%v spiderType:%v is not match", platform, spiderType)
		}
		return
	} else if platform == consts.Anjuke.Name() {
		switch spiderType {
		case consts.Area.Name():
			spider = anjuke.NewAreaSpider(city)
		case consts.ErShou.Name():
			spider = anjuke.NewErShouSpider(city)
		case consts.ChengJiao.Name():
			spider = anjuke.NewChengJiaoSpider(city)
		default:
			err = fmt.Errorf("current platform:%v spiderType:%v is not match", platform, spiderType)
		}
		return
	}
	err = fmt.Errorf("current platform:%v is not support", platform)
	return
}
