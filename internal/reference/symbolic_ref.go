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

var predefinedSymbolicRefTypes = []SymbolicRefType{
	Head,
	FetchHead,
	OrigHead,
	MergeHead,
}

type SymbolicRefParam struct {
	BaseDir       string
	Type          SymbolicRefType
	ReferencePath string
	Delete        bool
}

func SymbolicRef(param SymbolicRefParam) (string, error) {
	path := filepath.Join(param.BaseDir, string(param.Type))
	if param.Delete {
		if slices.Contains(predefinedSymbolicRefTypes, param.Type) {
			return "", fmt.Errorf("%s: refusing to delete predefined symbolic ref", param.Type)
		}
		return "", os.Remove(path)
	}

	if param.ReferencePath == "" {
		return readSymbolicRef(path)
	}

	if slices.Contains(predefinedSymbolicRefTypes, param.Type) && !strings.HasPrefix(param.ReferencePath, "refs/") {
		return "", fmt.Errorf("%s: refusing to point %s outside of refs/", param.ReferencePath, param.Type)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return "", err
	}

	return "", os.WriteFile(path, []byte(symbolicRefFormat(param.ReferencePath)), 0644)
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
