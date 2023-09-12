package file

import (
	"encoding/json"
	"github.com/byjayce/tig/internal/fn"
)

// Config is file.RegularFile implementation for configuration.
type Config struct {
	Core CoreConfig `json:"core"`
}

func (c Config) IsDir() bool {
	return false
}

func (c Config) Path() string {
	return "config"
}

func (c Config) Bytes() []byte {
	return fn.Must(json.Marshal(c))
}

func (c Config) String() string {
	return string(c.Bytes())
}

// CoreConfig is configuration for core.
type CoreConfig struct {
	Bare bool `json:"bare,omitempty"`
}
