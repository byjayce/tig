package porcelain

import (
	"errors"
	"fmt"
	"github.com/byjayce/tig/internal/reference"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type RevParseParam struct {
	TigDir string
	Target string // 미완성된 형태의 Hash, SymbolicRef, Ref 중 하나
}

// RevParse rev-parse 함수를 모방. 바라보고 있는 해시를 찾아주는 함수
func RevParse(param RevParseParam) (string, error) {
	_, err := os.Stat(filepath.Join(param.TigDir, param.Target))
	if err == nil {
		// 아무 에러가 없다면 -> SymbolicReference인 경우
		ref, err := reference.SymbolicRef(reference.SymbolicRefParam{
			TigDir: param.TigDir,
			Type:   reference.SymbolicRefType(param.Target),
		})
		if err != nil {
			return "", err
		}

		return RevParse(RevParseParam{
			TigDir: param.TigDir,
			Target: ref,
		})
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	// Reference인 경우
	var hash string
	err = filepath.WalkDir(filepath.Join(param.TigDir, "refs"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) != param.Target {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash = string(data)
		return fs.SkipAll
	})

	if hash != "" {
		return hash, nil
	}

	if err != nil {
		return "", err
	}

	//레퍼런스도 아니였다. 그렇다면 객체의 일부?
	if len(param.Target) < 4 {
		return "", fmt.Errorf("ambiguous argument '%s': unknown revision or path not in the working tree", param.Target)
	}

	err = filepath.WalkDir(filepath.Join(param.TigDir, "objects", param.Target[:2]), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasPrefix(filepath.Base(path), param.Target[2:]) {
			return nil
		}

		hash = param.Target[:2] + filepath.Base(path)
		return fs.SkipAll
	})

	if hash != "" {
		return hash, nil
	}

	return "", err
}
