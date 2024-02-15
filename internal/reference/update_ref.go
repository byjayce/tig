package reference

import (
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
	if param.Delete {
		return os.Remove(path)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(param.ObjectHash), 0644)
}
