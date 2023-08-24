package cfg

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DSN    string `mapstructure:"DB_DSN"`
	HTTP_HOST string `mapstructure:"HTTP_HOST"`
	HTTP_PORT string `mapstructure:"HTTP_PORT"`
}

var cfg Config

func Load(path string) error {
	readFromEnvFile(path)

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

func readFromEnvFile(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading env file, %s", err)
	}
}
