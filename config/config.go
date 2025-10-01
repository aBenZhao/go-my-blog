package config

import (
	"go-my-blog/pkg/logger"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// 配置读取
var Conf = new(AppConfig)

type AppConfig struct {
	Mysql  MysqlConfig  `mapstructure:"mysql"`
	Server ServerConfig `mapstructure:"server"`
	// 不在这里读取logConfig
	//Log    LogConfig    `mapstructure:"log"`
	// 不在这里读取GinConfig
	//Gin GinConfig `mapstructure:"gin"`
	JWT JWTConfig `mapstructure:"jwt"`
}

type MysqlConfig struct {
	DSN                 string `mapstructure:"dsn"`
	MaxOpenConns        int    `mapstructure:"max_open_conns"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetimeHour int    `mapstructure:"conn_max_lifetime_hour"`
	LogMode             bool   `mapstructure:"log_mode"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

//type LogConfig struct {
//	Level    string `mapstructure:"level"`
//	FilePath string `mapstructure:"file_path"`
//}

//type GinConfig struct {
//	Debug bool `mapstructure:"debug"`
//}

// JWTConfig JWT 配置结构体
type JWTConfig struct {
	Secret     string `mapstructure:"secret"`      // JWT 签名密钥
	ExpireHour int    `mapstructure:"expire_hour"` // Token 有效期（小时）
}

func Init() {
	// 读取配置文件
	viper.SetConfigFile("config/app.dev.yml")

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("读取数据库配置文件失败", zap.Error(err))
	}

	if err := viper.Unmarshal(Conf); err != nil {
		logger.Fatal("解析数据库配置到结构体失败", zap.Error(err))
	}

	// 验证配置
	validateMysqlConfig()
	logger.Info("配置初始化完成", zap.Any("mysql_config", Conf.Mysql))
}

func validateMysqlConfig() {
	if Conf.Mysql.DSN == "" {
		logger.Fatal("数据库DSN不能为空")
	}
	if Conf.Mysql.MaxOpenConns <= 0 {
		Conf.Mysql.MaxOpenConns = 50
		logger.Warn("数据库最大连接数不能小于等于0，已设置为默认值50")
	}
}

// GetConnMaxLifetime 辅助方法：将小时转为 time.Duration（供数据库连接池使用）
func (m *MysqlConfig) GetConnMaxLifetime() time.Duration {
	return time.Duration(m.ConnMaxLifetimeHour) * time.Hour
}
