package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	DirName = "objects"
)

type Type string

const (
	Blob   Type = "blob"
	Tree   Type = "tree"
	Commit Type = "commit"
	Tag    Type = "tag"
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

type Object struct {
	Type   Type
	Length int
	Data   []byte
}

func parseObject(tigDir string, objectHash string) (Object, error) {
	path := Key(objectHash).Path(tigDir)
	data, err := os.ReadFile(path)
	if err != nil {
		return Object{}, err
	}

	// 각 데이터에 맞게 Object 만들기
	var obj Object
	if err := unmarshalObject(data, &obj); err != nil {
		return Object{}, err
	}

	return obj, nil
}

func unmarshalObject(data []byte, o *Object) error {
	// 압축을 풀어준다.
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	decryption := new(bytes.Buffer)
	if _, err := decryption.ReadFrom(r); err != nil {
		return err
	}

	// Formatting HEADER\000BODY
	headerAndBody := bytes.Split(decryption.Bytes(), []byte{'\x00'})
	if len(headerAndBody) != 2 {
		return errors.New("invalid object format")
	}

	var (
		objectType   Type
		objectLength int
	)

	if _, err := fmt.Sscanf(string(headerAndBody[0]), "%s %d", &objectType, &objectLength); err != nil {
		return err
	}

	o.Type = objectType
	o.Length = objectLength
	o.Data = headerAndBody[1]
	return nil
}
