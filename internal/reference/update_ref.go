package reference

import (
	"errors"
	"os"
	"path/filepath"
)

type UpdateRefParam struct {
	BaseDir       string
	ReferencePath string
	ObjectHash    string
	Delete        bool
}

func UpdateRef(param UpdateRefParam) error {
	path := filepath.Join(param.BaseDir, param.ReferencePath)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	if param.Delete {
		err := os.Remove(path)
		if err == nil {
			return nil
		}

		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return err
	}

	return os.WriteFile(path, []byte(param.ObjectHash), 0644)
}
