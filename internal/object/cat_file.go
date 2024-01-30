package object

import (
	"errors"
	"fmt"
	"os"
)

type CatFileOperationType string

const (
	CatFileOperationTypePrettyPrint CatFileOperationType = "pretty-print" // CatFileOperationTypePrettyPrint pretty-print (-p)
	CatFileOperationTypeType        CatFileOperationType = "type"         // CatFileOperationTypeType type (-t)
	CatFileOperationTypeSize        CatFileOperationType = "size"         // CatFileOperationTypeSize size (-s)
	CatFileOperationTypeExist       CatFileOperationType = "exist"        // CatFileOperationTypeExist exist (-e)
)

type CatFileParam struct {
	BaseDir       string               // BaseDir tig 저장소 디렉토리
	OperationType CatFileOperationType // OperationType cat-file 옵션
	Hash          string               // Hash 객체 해시
}

func CatFile(param CatFileParam) (string, error) {
	objectPath := key(param.Hash).Path(param.BaseDir)
	b, err := os.ReadFile(objectPath)
	if err != nil {
		return "", err
	}

	v, err := parse(b)
	if err != nil {
		return "", err
	}

	switch param.OperationType {
	case CatFileOperationTypePrettyPrint:
		return string(v.Data), nil
	case CatFileOperationTypeType:
		return string(v.Type), nil
	case CatFileOperationTypeSize:
		return fmt.Sprintf("%d", len(v.Data)), nil
	case CatFileOperationTypeExist:
		return "", nil
	}

	return "", errors.New("invalid operation type")
}
