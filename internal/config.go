package internal

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ApiKey string `mapstructure:"NOTION_API_SECRET"`
}

func LoadConfig() (c Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return c, fmt.Errorf("cannot read config: %s", err)
	}

	if err := viper.Unmarshal(&c); err != nil {
		return c, fmt.Errorf("cannot unmarshal config: %s", err)
	}

	return
}
