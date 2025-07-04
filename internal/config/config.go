// internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
    CronSchedule   string `mapstructure:"cron_schedule"`
    GitHubToken    string `mapstructure:"github_token"`
    GitHubUsername string `mapstructure:"github_username"`
    GithubEmail    string `mapstructure:"github_email"`
}

var App AppConfig

func LoadConfig(path string) {
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	if err := viper.Unmarshal(&App); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

}

