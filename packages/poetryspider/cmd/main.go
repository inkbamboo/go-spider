package main

import (
	"fmt"
	"github.com/inkbamboo/go-spider/packages/poetryspider/internal/scripts"
	"github.com/labstack/gommon/color"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

const DateFullLayout = "2006-01-02 15:04:05"

func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:        "poetryspider",
			Aliases:     []string{"ps"},
			Usage:       "Run poetryspider",
			Description: "Run poetryspider",
			Action: func(c *cli.Context) error {
				fmt.Println("开始运行...")
				scripts.Init(c)
				time.Sleep(3 * time.Second)
				platform := c.String("platform")
				fmt.Println(fmt.Sprintf("%s%s%s%s%s%s",
					color.Bold(color.Green("platform:")), platform))
				scripts.RunPoetrySpider(platform)
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "Run test",
			Action: func(c *cli.Context) error {
				scripts.Init(c)
				scripts.RunTest()
				return nil
			},
		},
	}
}
func Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "conf",
			Value:       "./config",
			DefaultText: "./config",
			Usage:       "配置文件路径",
		},
		&cli.StringFlag{
			Name:        "env",
			DefaultText: "dev",
			Usage:       "执行环境 (开发环境dev、测试环境test、线上环境prod)",
		},
		&cli.StringFlag{
			Name:        "platform",
			DefaultText: "ke",
			Usage:       "指定平台",
		},
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "poetrypider"
	app.Usage = "诗词爬虫"
	app.Version = "1.0.0"
	app.Authors = []*cli.Author{{
		Name:  "inkbamboo",
		Email: "inkbamboo@icloud.com",
	}}
	app.Commands = Commands()
	app.Flags = Flags()
	_ = app.Run(os.Args)
	line := "==============================="
	fmt.Println(fmt.Sprintf("%s%s%s%s",
		color.White(line),
		color.Bold(color.Green("任务列表")),
		color.Bold(color.Yellow("["+time.Now().Format(DateFullLayout))+"]"),
		color.White(line)))
	fmt.Println(color.Bold(color.White("包含以下任务:")))
	for key, command := range app.Commands {
		fmt.Println(fmt.Sprintf("任务%d：%s %s %s", key+1, command.Name, command.Usage, command.Description))
	}
}
