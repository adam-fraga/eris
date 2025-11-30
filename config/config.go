package config

import (
	"fmt"

	v "github.com/spf13/viper"
)

type Config struct {
	Url          string `mapstructure:"url" env:"URL" json:"url"`
	Model        string `mapstructure:"model" env:"MODEL" json:"model"`
	SystemPrompt string `mapstructure:"systemPrompt" env:"SYSTEM_PROMPT" json:"systemPrompt"`
}

func LoadConfig() (*Config, error) {

	v.SetDefault("url", "http://localhost:11434/api/chat")
	v.SetDefault("model", "gpt-oss:20b")
	v.SetDefault("systemPrompt", "You are an expert in computer programming Please make friendly answer for the noobs. Add source code examples if you can.")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
