// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	CronSchedule string `mapstructure:"cron_schedule"`
	GitHubToken  string `mapstructure:"github_token"`
}

var App AppConfig

func LoadConfig(path string) {
	viper.SetConfigFile(path)
	viper.SetEnvPrefix("SYNC")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&App); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// fallback to env variables
	if token := os.Getenv("SYNC_GITHUB_TOKEN"); token != "" {
		App.GitHubToken = token
	}
}

