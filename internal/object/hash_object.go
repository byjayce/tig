package object

import (
	"os"
)

type HashObjectParam struct {
	DryRun bool   // DryRun -w 옵션이 켜져있는 경우 DryRun = false, vice versa
	TigDir string // TigDir tig 저장소 디렉토리 경로
	Type   Type   // Type 객체의 타입 종류
	Data   []byte // Data 객체 내용
}

func HashObject(param HashObjectParam) (string, error) {
	// 키 만들기
	key := newKey(param.Type, param.Data)

	// DryRun 여기서 퀵 리턴
	if param.DryRun {
		return string(key), nil
	}

	// 내용 만들기 (zlib 압축)
	val, err := zlibCompress(newContent(param.Type, param.Data))
	if err != nil {
		return "", err
	}

	// 필요한 디렉토리 만들기 (Key 가장 앞 2자리 디렉토리까지)
	if err := os.MkdirAll(key.Dir(param.TigDir), os.ModePerm); err != nil {
		return "", err
	}

	// 필요한 객체 파일 생성 (쓰기)
	if err := os.WriteFile(key.Path(param.TigDir), val, os.ModePerm); err != nil {
		return "", err
	}

	return string(key), nil
}
