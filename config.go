package main

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	interval int
)

func readConfig() {

	viper.SetConfigName("config.txt") // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")          // optionally look for config in the working directory

	// Handle errors reading the config file
	if err := viper.ReadInConfig(); err != nil {
		logger.Error().Str("error", err.Error()).Msg("read config.txt error")
		fmt.Printf("解析配置文件出错，错误原因: %s\n", err.Error())
		waitExit()
	}

	interval = viper.GetInt("interval")
	logger.Info().Int("interval", interval).Msg("config")

}
