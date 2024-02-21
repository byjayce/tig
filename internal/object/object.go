package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"path/filepath"
)

const (
	DirName = "objects"
)

type Type string

const (
	Blob Type = "blob"
)

type Key string

func (k Key) Dir(tigDir string) string {
	return filepath.Join(tigDir, DirName, string(k[:2]))
}

func (k Key) Path(tigDir string) string {
	return filepath.Join(k.Dir(tigDir), string(k[2:]))
}

func newKey(t Type, data []byte) Key {
	// [type, length, \0, data] -> SHA-1 Hash
	str := newContent(t, data)

	// SHA-1 Hash
	h := sha1.New()
	h.Write(str)

	// 16진수 문자열로 변환
	return Key(fmt.Sprintf("%x", h.Sum(nil)))
}

func newContent(t Type, data []byte) []byte {
	// [Type] [Length of Type]\0[Data]
	return []byte(fmt.Sprintf("%s %d\000%s", t, len(data), data))
}

func zlibCompress(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)

	// zlib writer 생성, Buffer 쓰기
	w := zlib.NewWriter(b)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}

	err := w.Close()
	return b.Bytes(), err
}
