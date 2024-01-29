package workingcopy

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	configFileName = "config"
	headFileName   = "HEAD"
	objectsDirName = "objects"
	refsDirName    = "refs"
)

type InitCoreConfig struct {
	Bare bool `yaml:"bare"` // Bare 옵션. Bare 모드인지 아닌지 여부를 지정한다.
}

type InitConfig struct {
	Core InitCoreConfig `yaml:"core"` // Core 옵션.
}

// InitParam
// 이 구조체는 Init() 함수의 파라미터로 사용된다.
type InitParam struct {
	WorkingCopyPath string     // 작업 공간의 경로
	Config          InitConfig // Init() 함수가 `config` 파일을 만들 때 지정할 내용들
}

// Init
// 이 함수는 Git의 작업 공간을 초기화한다.
func Init(param InitParam) error {
	// TODO: 설정에 따라 Base가 바뀌는 상황 추가하기
	base := param.WorkingCopyPath

	if err := createConfigFile(base, param.Config); err != nil {
		return err
	}

	if err := createHeadFile(base); err != nil {
		return err
	}

	if err := createObjectsDir(base); err != nil {
		return err
	}

	return createRefsDir(base)
}

func createConfigFile(base string, param InitConfig) error {
	buf, err := yaml.Marshal(param)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(base, configFileName), buf, 0644)
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
