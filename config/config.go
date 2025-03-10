package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	ModelName      string        `json:"model_name"`
	OutPath        string        `json:"out_path"`
	LogsPath       string        `json:"logs_path"`
	LogLines       int           `json:"log_lines"`
	ParagraphChars int           `json:"paragraph_chars"`
	TotalChars     int           `json:"total_chars"`
	SleepTime      time.Duration `json:"sleep_time"`
	SystemPrompt   []string      `json:"system_prompt"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) Render() {
	for i, line := range c.SystemPrompt {
		line = strings.ReplaceAll(line, "{paragraph_chars}", fmt.Sprintf("%d", c.ParagraphChars))
		line = strings.ReplaceAll(line, "{total_chars}", fmt.Sprintf("%d", c.TotalChars))
		c.SystemPrompt[i] = line
	}
}
