package cfg

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DSN string `mapstructure:"DB_DSN"`

	API_HOST string `mapstructure:"API_HOST"`
	API_PORT string `mapstructure:"API_PORT"`

	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_PORT string `mapstructure:"REDIS_PORT"`
}

var cfg Config

func Load(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading env file, %s", err)
	}

	viper.AutomaticEnv()

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to unmarshall into struct, %v", err)
	}

	return nil
}

func Get() *Config {
	return &cfg
}
