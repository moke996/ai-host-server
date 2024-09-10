package global

import (
	"ai-dating/model"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config model.Config

func LoadConfig() {
	confPath := "./config.yaml"
	v := viper.New()
	v.SetConfigFile(confPath) // 配置文件的路径

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}
	// monitor the changes in the config file
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.ReadInConfig(); err != nil {
			panic(err)
		}

		if err := v.Unmarshal(&Config); err != nil {
			panic(err)
		}
	})
}
