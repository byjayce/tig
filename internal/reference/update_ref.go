package reference

import (
	"errors"
	"os"
	"path/filepath"
)

type UpdateRefParam struct {
	TigDir        string
	ReferencePath string // tig 디렉토리 기준 레퍼런스 파일의 위치
	ObjectHash    string
	Delete        bool
}

func UpdateRef(param UpdateRefParam) error {
	refPath := filepath.Join(param.TigDir, param.ReferencePath)
	if err := os.MkdirAll(filepath.Dir(refPath), os.ModePerm); err != nil {
		return err
	}

	if !param.Delete {
		return os.WriteFile(refPath, []byte(param.ObjectHash), os.ModePerm)
	}

	err := os.Remove(refPath)
	if err == nil {
		return nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}
