package config

import (
	"github.com/spf13/viper"
)

func GetConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		Marvel: Marvel{
			PublicAPIKey:  viper.GetString("MARVEL_API_PUBLIC"),
			PrivateAPIKey: viper.GetString("MARVEL_API_PRIVATE"),
		},
	}
	return config, nil
}
