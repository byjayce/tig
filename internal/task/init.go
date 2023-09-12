package task

import (
	"github.com/byjayce/tig/internal/domain"
	"github.com/byjayce/tig/internal/file"
)

// InitParam is parameter model for init task.
type InitParam struct {
	// Bare is flag for bare repository.
	// If true, repository is initialized as bare repository.
	// If false, repository is initialized as normal repository.
	Bare bool

	// InitParam.InitialBranch is initial branch name.
	InitialBranch string

	// InitParam.Repository is Directory implementation for current working directory.
	Repository domain.TigRepository
}

func Init(param InitParam) error {
	// if current directory is already initialized, do nothing.
	if param.Repository.IsInitialized() {
		return nil
	}

	return param.Repository.Init(
		file.Config{
			Core: file.CoreConfig{
				Bare: param.Bare,
			},
		},
		file.SymbolicReference{
			Point:    param.InitialBranch,
			FilePath: []string{file.SymbolicHead},
		},
	)
}
