package file

import (
	"fmt"
	"path/filepath"
)

const (
	SymbolicHead = "HEAD"
)

type SymbolicReference struct {
	Point    string
	FilePath []string
}

func (s SymbolicReference) Path() string {
	return filepath.Join(s.FilePath...)
}

func (s SymbolicReference) Bytes() []byte {
	return []byte(s.String())
}

func (s SymbolicReference) String() string {
	return fmt.Sprintf("ref: %s\n", s.Point)
}
