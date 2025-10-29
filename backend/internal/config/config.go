package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Worker   WorkerConfig   `mapstructure:"worker"`
	Whisper  WhisperConfig  `mapstructure:"whisper"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	UploadPath  string `mapstructure:"upload_path"`
	MaxFileSize int64  `mapstructure:"max_file_size"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Name     string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type OpenAIConfig struct {
	APIBase     string `mapstructure:"api_base"`
	APIKey      string `mapstructure:"api_key"`
	ModelVision string `mapstructure:"model_vision"`
	ModelChat   string `mapstructure:"model_chat"`
}

type JWTConfig struct {
	Secret         string `mapstructure:"secret"`
	Expiry         string `mapstructure:"expiry"`
	RefreshSecret  string `mapstructure:"refresh_secret"`
	RefreshExpiry  string `mapstructure:"refresh_expiry"`
}

type WorkerConfig struct {
	Concurrency int `mapstructure:"concurrency"`
}

type WhisperConfig struct {
	ServiceURL string `mapstructure:"service_url"`
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.upload_path", "./uploads")
	viper.SetDefault("server.max_file_size", 524288000) // 500MB

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.name", "opinion_monitor")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "password")

	viper.SetDefault("openai.api_base", "https://api.openai.com/v1")
	viper.SetDefault("openai.model_vision", "gpt-4o")
	viper.SetDefault("openai.model_chat", "gpt-4o")

	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expiry", "15m")
	viper.SetDefault("jwt.refresh_secret", "your-refresh-secret-key-change-in-production")
	viper.SetDefault("jwt.refresh_expiry", "168h") // 7 days

	viper.SetDefault("worker.concurrency", 5)

	viper.SetDefault("whisper.service_url", "http://localhost:5000")

	// Allow environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
