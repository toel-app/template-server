package utils

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"github.com/toel-app/common-utils/logger"
)

func init() {
	env := os.Getenv("env")

	if len(env) == 0 {
		env = "dev"
	}

	defaultConfig := fmt.Sprintf("config.%s", env)

	viper.SetConfigName(defaultConfig)
	viper.SetConfigType("yml")
	viper.AddConfigPath("manifests")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("config not properly configured", nil)
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
}
