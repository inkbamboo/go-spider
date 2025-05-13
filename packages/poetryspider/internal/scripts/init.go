package scripts

import (
	"fmt"
	"github.com/inkbamboo/ares"
	"github.com/urfave/cli/v2"
)

func Init(c *cli.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Error: %v\n", r)
		}
	}()
	platform := c.String("platform")
	ares.InitConfigWithPath(c.String("env"), c.String("conf"))
	for _, database := range ares.GetBaseConfig().Databases {
		database.Alias = fmt.Sprintf("%s_poetry", platform)
		database.DbName = fmt.Sprintf("%s_poetry", platform)
	}
	ares.GetConfig().Set("env", c.String("env"))
	ares.NewAres()
}
