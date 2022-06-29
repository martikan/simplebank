package util

import (
	"time"

	"github.com/spf13/viper"
)

var (
	ConfigUtils configUtilsInterface = &configUtils{}
)

type configUtils struct{}

type configUtilsInterface interface {
	LoadConfig(string) (Config, error)
}

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	DbDriver            string        `mapstructure:"DATASOURCE_DRIVER"`
	DbUrl               string        `mapstructure:"DATASOURCE_URL"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

// LoadConfig reads the configuration from file or environment variables.
func (c *configUtils) LoadConfig(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.SetConfigName("application")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
