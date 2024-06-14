package pkg

import (
	"fmt"
	"github.com/spf13/viper"
)

var (
	v *viper.Viper
)

func InitConfig(env string) {
	InitConfigWithPath(env, "./config/")
}
func InitConfigWithPath(env string, configPath string) {
	fmt.Println(fmt.Sprintf("配置文件路径: %s", configPath))
	fmt.Println(fmt.Sprintf("执行环境: %s", env))
	v = viper.New()
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println(fmt.Sprintf("Viper ReadInConfig err:%s\n", err))
		panic(err)
	}
	v.Set("env", env)
}
func GetConfig() *viper.Viper {
	if v == nil {
		panic("Please init Config")
	}
	return v
}

func GetEnv() string {
	return v.GetString("env")
}
