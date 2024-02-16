package porcelain

import (
	"errors"
	"github.com/byjayce/tig/internal/config"
	"os"
	"path/filepath"
)

const baseDir = ".git"

type Tig struct {
	config  config.Config
	baseDir string
}

func NewTig() (*Tig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gitPath := filepath.Join(wd, baseDir)

	stat, err := os.Stat(gitPath)
	if err == nil {
		if stat.IsDir() {
			cfg, err := config.ReadConfigFile(gitPath)
			if err != nil {
				return nil, err
			}

			return &Tig{
				config:  cfg,
				baseDir: baseDir,
			}, nil
		}
	}

	if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	cfg, err := config.ReadConfigFile(wd)
	if err != nil {
		return nil, err
	}

	if !cfg.Core.Bare {
		return nil, errors.New("not a git repository")
	}

	return &Tig{
		config:  cfg,
		baseDir: wd,
	}, nil
}
