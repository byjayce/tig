package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	configFileName = "config"
)

type Core struct {
	Bare bool // Bare 옵션. Bare 모드인지 아닌지 여부를 지정한다.
}

type User struct {
	Name  string // 사용자 이름
	Email string // 사용자 이메일
}

type Config struct {
	Core Core // Core 옵션.
	User User // User 옵션.
}

func CreateConfigFile(base string, param Config) error {
	buf, err := json.Marshal(param)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(base, configFileName), buf, 0644)
}

func ReadConfigFile(base string) (Config, error) {
	buf, err := os.ReadFile(filepath.Join(base, configFileName))
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(buf, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
