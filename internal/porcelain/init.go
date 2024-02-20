package porcelain

import (
	"github.com/byjayce/tig/internal/config"
	"os"
	"path/filepath"
)

const (
	headFileName   = "HEAD"
	objectsDirName = "objects"
	refsDirName    = "refs"
)

// InitParam
// 이 구조체는 Init() 함수의 파라미터로 사용된다.
type InitParam struct {
	WorkingCopyDir string        // 작업 공간의 경로
	Config         config.Config // Init() 함수가 `config` 파일을 만들 때 지정할 내용들
}

// Init
// 이 함수는 Git의 작업 공간을 초기화한다.
func Init(param InitParam) error {
	tigDir := param.WorkingCopyDir
	if !param.Config.Core.Bare {
		tigDir = filepath.Join(tigDir, baseDir)
		if err := os.MkdirAll(tigDir, 0755); err != nil {
			return err
		}
	}

	if err := config.CreateConfigFile(tigDir, param.Config); err != nil {
		return err
	}

	if err := createHeadFile(tigDir); err != nil {
		return err
	}

	if err := createObjectsDir(tigDir); err != nil {
		return err
	}

	return createRefsDir(tigDir)
}

func createHeadFile(base string) error {
	// TODO: Symbolic Ref를 저장하는 포맷을 만들도록 수정하기
	buf := []byte("ref: refs/heads/main")
	return os.WriteFile(filepath.Join(base, headFileName), buf, 0644)
}

func createObjectsDir(base string) error {
	return os.MkdirAll(filepath.Join(base, objectsDirName), 0755)
}

func createRefsDir(base string) error {
	return os.MkdirAll(filepath.Join(base, refsDirName), 0755)
}
