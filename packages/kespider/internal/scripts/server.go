package scripts

import (
	"github.com/inkbamboo/go-spider/packages/kespider/internal/spider"
)

func RunAreaSpider() {
	spider.GetAreaSpider().Start()
	select {}
}
func RunErShouSpider() {
	spider.GetErShouSpider().Start()
	select {}
}
