package object

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
)

const (
	objectsDirName = "objects"
)

type key string

func newKey(t Type, data []byte) key {
	str := content(t, data)
	h := sha1.New()
	h.Write([]byte(str))
	return key(fmt.Sprintf("%x", h.Sum(nil)))
}

func (k key) Dir(base string) string {
	return filepath.Join(base, objectsDirName, string(k[:2]))
}

func (k key) Path(base string) string {
	return filepath.Join(base, objectsDirName, string(k[:2]), string(k[2:]))
}

func (k key) String() string {
	return string(k)
}
