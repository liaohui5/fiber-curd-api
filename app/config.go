package app

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config 全局配置实例
var Config *viper.Viper

// InitConfig 初始化配置
func InitConfig() {
	Config = viper.New()

	// 设置配置文件路径和类型
	Config.SetConfigType("toml")
	Config.SetConfigName("config")
	Config.AddConfigPath(".")

	// 读取配置
	if err := Config.ReadInConfig(); err != nil {
		log.Fatalf("读取基础配置失败: %v \n", err)
	}

	// 加载环境配置文件
	isDev := strings.ToLower(os.Getenv("IS_DEV"))
	var envConf string
	if isDev == "true" || isDev == "yes" || isDev == "1" {
		envConf = "config.dev"
		fmt.Printf("[DEV] %s.toml loaded \n", envConf)
	} else {
		envConf = "config.prod"
		fmt.Printf("[PROD] %s.toml loaded \n", envConf)
	}

	// 合并配置
	Config.SetConfigName(envConf)
	if err := Config.MergeInConfig(); err != nil {
		fmt.Printf("[WARN]: %s not found\n", envConf)
	}
}
