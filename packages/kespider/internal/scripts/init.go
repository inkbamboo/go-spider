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
	ares.InitConfigWithPath(c.String("env"), c.String("conf"))
	ares.GetConfig().Set("env", c.String("env"))
	ares.NewAres()
}
