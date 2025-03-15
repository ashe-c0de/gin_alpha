package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Dbname   string `mapstructure:"dbname"`
	Port     string `mapstructure:"port"`
	TimeZone string `mapstructure:"timezone"`
	SslMode  string `mapstructure:"sslmode"`
}

// LoadConfig 读取 YAML 配置
func LoadConfig() (*DatabaseConfig, error) {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config DatabaseConfig
	if err := viper.UnmarshalKey("database", &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal database config: %w", err)
	}
	return &config, nil
}
