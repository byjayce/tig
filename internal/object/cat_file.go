package object

import (
	"errors"
	"fmt"
)

type CatFileOperationType string

const (
	CatFileOperationTypePrettyPrint CatFileOperationType = "pretty-print" // -p 옵션인 경우
	CatFileOperationTypeType        CatFileOperationType = "type"         // -t 옵션인 경우
	CatFileOperationTypeSize        CatFileOperationType = "size"         // -s 옵션인 경우
	CatFileOperationTypeExist       CatFileOperationType = "exist"        // -e 옵션인 경우
)

type CatFileParam struct {
	TigDir        string
	OperationType CatFileOperationType
	ObjectHash    string
}

func CatFile(param CatFileParam) (string, error) {
	// 객체를 디렉토리에서 읽은 다음 내용을 가져와야 함
	obj, err := parseObject(param.TigDir, param.ObjectHash)
	if err != nil {
		return "", err
	}

	// 목적에 맞는 리턴을 해줘야 함
	switch param.OperationType {
	case CatFileOperationTypePrettyPrint:
		return string(obj.Data), nil
	case CatFileOperationTypeType:
		return string(obj.Type), nil
	case CatFileOperationTypeSize:
		return fmt.Sprintf("%d", obj.Length), nil
	case CatFileOperationTypeExist:
		return "", nil
	}

	// 정해진 목적 (Operation Type)이 지원되지 않는 경우 에러를 리턴
	return "", errors.New("implement me")
}
