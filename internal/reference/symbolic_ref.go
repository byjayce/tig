package reference

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type SymbolicRefType string

const (
	Head      SymbolicRefType = "HEAD"
	FetchHead SymbolicRefType = "FETCH_HEAD"
	OrigHead  SymbolicRefType = "ORIG_HEAD"
	MergeHead SymbolicRefType = "MERGE_HEAD"
)

var predefinedTypes = []SymbolicRefType{
	Head,
	FetchHead,
	OrigHead,
	MergeHead,
}

type SymbolicRefParam struct {
	TigDir        string
	Type          SymbolicRefType
	ReferencePath string
	Delete        bool
}

func SymbolicRef(param SymbolicRefParam) (string, error) {
	path := filepath.Join(param.TigDir, string(param.Type))

	// 삭제 옵션이 켜진 경우
	if param.Delete {
		// 미리 정의된 타입은 지워질 수 없음.
		if slices.Contains(predefinedTypes, param.Type) {
			return "", fmt.Errorf("%s: refusing to delete predefined symbolic ref", param.Type)
		}

		return "", os.RemoveAll(path)
	}

	if param.ReferencePath == "" {
		// 단순하게 Symbolic Reference를 읽는 동작
		return readSymbolicRef(path)
	}

	if slices.Contains(predefinedTypes, param.Type) && !strings.HasPrefix(param.ReferencePath, "refs/") {
		return "", fmt.Errorf("%s: refusing to point %s outside of refs/", param.ReferencePath, param.Type)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}

	return "", os.WriteFile(path, []byte(symbolicRefFormat(param.ReferencePath)), os.ModePerm)
}

func readSymbolicRef(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return strings.TrimPrefix(string(data), "ref: "), nil
}

func symbolicRefFormat(refPath string) string {
	return fmt.Sprintf("ref: %s", refPath)
}
