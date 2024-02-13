package object

import (
	"os"
)

type HashObjectParam struct {
	DryRun  bool
	BaseDir string // BaseDir tig 저장소 디렉토리 경로
	Type    Type   // Type 객체 타입
	Data    []byte // Data 객체 내용
}

// HashObject
// 객체를 `objects` 디렉토리 아래 규칙에 맞게 저장한다.
func HashObject(param HashObjectParam) (string, error) {
	k := newKey(param.Type, param.Data)
	if param.DryRun {
		return k.String(), nil
	}

	val, err := zlibCompress(content(param.Type, param.Data))
	if len(val) == 0 {
		return "", err
	}

	if err := os.MkdirAll(k.Dir(param.BaseDir), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(k.Path(param.BaseDir), val, 0644); err != nil {
		return "", err
	}

	// Close 실패한 경우 에러가 있을 수 있음
	return k.String(), err
}
