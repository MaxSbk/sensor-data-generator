package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
)

func LoadConfig(filePath string) (Config, error) {
	absPath, _ := filepath.Abs(filePath)
	var cfg Config
	err := cleanenv.ReadConfig(absPath, &cfg)
	return cfg, err
}
