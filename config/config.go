package config

import (
	"fmt"

	v "github.com/spf13/viper"
)

type Config struct {
	Url           string `mapstructure:"url" env:"URL" json:"url"`
	ThinkingModel string `mapstructure:"thinkingModel" env:"THINKING_MODEL" json:"thinkingModel"`
	CodingModel   string `mapstructure:"codingModel" env:"CODING_MODEL" json:"codingModel"`
	SystemPrompt  string `mapstructure:"systemPrompt" env:"SYSTEM_PROMPT" json:"systemPrompt"`
}

func LoadConfig() (*Config, error) {

	v.SetDefault("url", "http://localhost:11434/api/chat")
	v.SetDefault("codingModel", "qwen3-coder:30b")
	v.SetDefault("thinkingModel", "qwen3-vl:30b")
	v.SetDefault("systemPrompt", `
		You are an expert AI assistant specialized in computer programming and automation. 
		You help users automate daily tasks using local tools and APIs.  
		You provide clear, friendly explanations, examples, and working code snippets when possible. 
		Always consider the available local functions (tools) and suggest using them when relevant. 
		Focus on practical, actionable answers that a non-expert can follow. 
		Be concise, but thorough enough for the user to implement solutions quickly.`)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}
