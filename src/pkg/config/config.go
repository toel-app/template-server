package config

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/toel-app/template-server/src/pkg/logger"

	"os"
	"sync"
)

type Config struct {
	Database struct {
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
		Name string `mapstructure:"name"`
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"database"`
	Keycloak struct {
		Host          string `mapstructure:"host"`
		Realm         string `mapstructure:"realm"`
		AdminUser     string `mapstructure:"adminUser"`
		AdminPassword string `mapstructure:"adminPassword"`
		ClientSecret  string `mapstructure:"clientSecret"`
		ClientId      string `mapstructure:"clientId"`
	} `mapstructure:"keycloak"`
}

var (
	cfg  Config
	once sync.Once
)

func NewConfig(logger logger.Logger) Config {
	once.Do(func() {
		env, ok := os.LookupEnv("env")
		if !ok || len(env) == 0 {
			env = "dev"
		}

		defaultConfig := fmt.Sprintf("config.%s", env)

		viper.SetConfigName(defaultConfig)
		viper.SetConfigType("yml")
		viper.AddConfigPath("manifests")
		viper.AddConfigPath(".")

		if err := viper.ReadInConfig(); err != nil {
			logger.Error("config not properly configured", nil)
			panic(fmt.Errorf("Fatal error config file: %w \n", err))
		}

		var config Config

		opt := func(config *mapstructure.DecoderConfig) {
			config.ErrorUnused = true
			config.ErrorUnset = true
		}
		if err := viper.Unmarshal(&config, opt); err != nil {
			panic(err)
		}

		cfg = config
	})

	return cfg
}
