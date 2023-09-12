package repository

import (
	"errors"
	"github.com/byjayce/tig/internal/file"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type Repository struct {
	path   string
	dir    fs.FS
	logger *log.Logger
}

func New(logger *log.Logger, wd string) *Repository {

	return &Repository{
		logger: logger,
		dir:    os.DirFS(wd),
		path:   wd,
	}
}

func (r *Repository) IsInitialized() bool {
	if _, err := fs.Stat(r.dir, file.SymbolicHead); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}

		r.logger.Fatal(err)
	}

	return true
}

func (r *Repository) Init(config file.Config, head file.SymbolicReference) error {
	var (
		refsHeads = filepath.Join(r.path, "refs", "heads")
		objects   = filepath.Join(r.path, "objects")
	)

	if err := os.MkdirAll(refsHeads, 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(objects, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(file.SymbolicHead, head.Bytes(), 0644); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(r.path, "config"), config.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}
