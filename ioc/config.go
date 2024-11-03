package ioc

import (
	"github.com/WeiXinao/hi-wiki/config"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitViper() *config.AppConfig {
	configFile := pflag.String("config", "config/dev.yaml", "配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	var AppCfg config.AppConfig
	err = viper.Unmarshal(&AppCfg)
	if err != nil {
		panic(err)
	}
	return &AppCfg
}
