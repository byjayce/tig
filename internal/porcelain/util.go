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

// resolveRef 심볼릭 레퍼런스 혹은 레퍼런스를 사용해 객체 해시를 찾는다. (rev-parse 유사 동작)
func resolveRef(base, ref string) (string, error) {
	_, err := os.Stat(filepath.Join(base, ref))
	if err == nil {
		// Symbolic Reference인 경우
		ref, err = reference.SymbolicRef(reference.SymbolicRefParam{
			BaseDir: baseDir,
			Type:    reference.SymbolicRefType(ref),
		})
		if err != nil {
			return "", err
		}
		return resolveRef(base, ref)
	}

	if !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	// Reference인 경우, refs 디렉토리 아래를 순회하고 파일이름이 같은 걸 찾음
	var hash string
	err = filepath.WalkDir(filepath.Join(base, "refs"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) != ref {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash = string(data)
		return fs.SkipAll
	})

	if errors.Is(err, fs.SkipAll) {
		return hash, nil
	}

	if err == nil {
		// 레퍼런스가 아님. object hash의 일부인지 확인
		if len(ref) < 4 {
			return "", fmt.Errorf("ambiguous argument '%s': unknown revision or path not in the working tree", ref)
		}

		err = filepath.WalkDir(filepath.Join(base, "objects", ref[:2]), func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if !strings.HasPrefix(filepath.Base(path), ref[2:]) {
				return nil
			}

			hash = ref
			return fs.SkipAll
		})

		if errors.Is(err, fs.SkipAll) {
			return hash, nil
		}
	}

	return hash, err
}
