package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type User struct {
	Email string
	Name  string
}

type Core struct {
	Bare bool
}

type Config struct {
	Core Core
	User User
}

const configFileName = "config"

func CreateConfigFile(tigDir string, cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(tigDir, configFileName), data, os.ModePerm)
}

func ReadConfigFile(tigDir string) (Config, error) {
	// 파일을 읽음
	data, err := os.ReadFile(filepath.Join(tigDir, configFileName))
	if err != nil {
		return Config{}, err
	}

	// Json -> Config 자료구조로 옮김
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, nil
	}

	return cfg, nil
}
