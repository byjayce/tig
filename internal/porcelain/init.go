package porcelain

import (
	"github.com/byjayce/tig/internal/config"
	"github.com/byjayce/tig/internal/object"
	"os"
	"path/filepath"
)

const (
	baseDir        = ".tig"
	headFileName   = "HEAD"
	objectsDirName = "objects"
	refsDirName    = "refs"
)

type InitParam struct {
	WorkingCopyDir string
	Config         config.Config
}

func Init(param InitParam) error {
	// config 파일 만들기
	tigDir := param.WorkingCopyDir
	if !param.Config.Core.Bare {
		// Bare 옵션이 꺼진 경우
		tigDir = filepath.Join(tigDir, baseDir)
		if err := os.MkdirAll(tigDir, os.ModePerm); err != nil {
			return err
		}
	}

	if err := config.CreateConfigFile(tigDir, param.Config); err != nil {
		return err
	}

	// HEAD 파일 만들기
	if err := createHeadFile(tigDir); err != nil {
		return err
	}

	// objects 디렉토리 만들기
	if err := os.MkdirAll(filepath.Join(tigDir, object.DirName), os.ModePerm); err != nil {
		return err
	}

	// refs 디렉토리 만들기
	return os.MkdirAll(filepath.Join(tigDir, refsDirName), os.ModePerm)
}

func createHeadFile(tigDir string) error {
	// TODO: 직접 쓰기 보단, Internal 명령어로 대체할 수 있으면 사용하기
	data := []byte("ref: refs/heads/main")
	return os.WriteFile(filepath.Join(tigDir, headFileName), data, os.ModePerm)
}
