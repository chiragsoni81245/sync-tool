// internal/config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
    CronSchedule                   string `mapstructure:"cron_schedule"`
    DatabaseFilepath               string `mapstructure:"database_filepath"`
    GitHubToken                    string `mapstructure:"github_token"`
    GitHubUsername                 string `mapstructure:"github_username"`
    GithubEmail                    string `mapstructure:"github_email"`
    GoogleDriveCredentialsFilepath string `mapstructure:"google_drive_credentials_filepath"`
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

