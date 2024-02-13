package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"os"
)

// Type 객체 타입
type Type string

const (
	BlobType   Type = "blob"   // BlobType Blob 객체 타입
	TreeType   Type = "tree"   // TreeType Tree 객체 타입
	CommitType Type = "commit" // CommitType Commit 객체 타입
)

type value struct {
	Type Type
	Data []byte
}

func (v value) Bytes() ([]byte, error) {
	return zlibCompress(content(v.Type, v.Data))
}

func parseObject(baseDir, objectKey string) (value, error) {
	objectPath := key(objectKey).Path(baseDir)
	data, err := os.ReadFile(objectPath)
	if err != nil {
		return value{}, err
	}

	return parse(data)
}

func parse(data []byte) (value, error) {
	b, err := parseRaw(data)
	if err != nil {
		return value{}, err
	}
	return parseFormat(b)
}

func parseRaw(data []byte) ([]byte, error) {
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	if _, err := b.ReadFrom(r); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func parseFormat(data []byte) (value, error) {
	var (
		objectType    Type
		objectSize    int
		objectContent string
	)

	headerAndBody := bytes.Split(data, []byte{'\x00'})
	if len(headerAndBody) != 2 {
		return value{}, fmt.Errorf("invalid object format")
	}

	if _, err := fmt.Sscanf(string(headerAndBody[0]), "%s %d", &objectType, &objectSize); err != nil {
		return value{}, err
	}
	objectContent = string(headerAndBody[1])

	return value{
		Type: objectType,
		Data: []byte(objectContent),
	}, nil
}

func zlibCompress(c string) ([]byte, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	if _, err := w.Write([]byte(c)); err != nil {
		return nil, err
	}

	err := w.Close()
	return b.Bytes(), err
}
