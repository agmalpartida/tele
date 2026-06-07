package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type TelegramConfig struct {
	APIID       int    `mapstructure:"api_id"`
	APIHash     string `mapstructure:"api_hash"`
	SessionFile string `mapstructure:"session_file"`
}

type UIConfig struct {
	Theme        string `mapstructure:"theme"`
	DateFormat   string `mapstructure:"date_format"`
	HistoryLimit int    `mapstructure:"history_limit"`
}

type PhotosConfig struct {
	EagerFullQuality bool   `mapstructure:"eager_full_quality"`
	Mode             string `mapstructure:"mode"` // auto | kitty | blocks
}

type Config struct {
	Telegram TelegramConfig `mapstructure:"telegram"`
	UI       UIConfig       `mapstructure:"ui"`
	Photos   PhotosConfig   `mapstructure:"photos"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	setDefaults(v)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	if cfg.Telegram.SessionFile == "" {
		cfg.Telegram.SessionFile = filepath.Join(filepath.Dir(path), "session.json")
	} else if strings.HasPrefix(cfg.Telegram.SessionFile, "~/") {
		home, _ := os.UserHomeDir()
		cfg.Telegram.SessionFile = filepath.Join(home, cfg.Telegram.SessionFile[2:])
	}
	return &cfg, nil
}
