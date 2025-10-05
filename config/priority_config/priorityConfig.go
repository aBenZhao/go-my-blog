package priority_config

import (
	"log"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 配置读取
var PriorityConf = new(PriorityConfig)

type PriorityConfig struct {
	Log LogConfig `mapstructure:"log"`
	Gin GinConfig `mapstructure:"gin"`
}

type LogConfig struct {
	Level    string `mapstructure:"level"`
	FilePath string `mapstructure:"file_path"`
}

type GinConfig struct {
	Debug bool `mapstructure:"debug"`
}

func PriorityConfInit() {
	// 读取配置文件
	viper.SetConfigFile("config/app.dev.yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("读取数据库配置文件失败", zap.Error(err))
		//logger.Fatal("读取数据库配置文件失败", zap.Error(err))
	}

	if err := viper.Unmarshal(PriorityConf); err != nil {
		log.Fatal("解析数据库配置到结构体失败", zap.Error(err))
	}
}
