package scripts

import (
	"github.com/inkbamboo/go-spider/packages/pkg"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	pkg.InitConfig(c.String("env"))
}
