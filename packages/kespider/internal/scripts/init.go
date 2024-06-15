package scripts

import (
	"github.com/inkbamboo/ares"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	ares.InitConfigWithPath(c.String("env"), c.String("conf"))
	ares.GetConfig().Set("env", c.String("env"))
}
