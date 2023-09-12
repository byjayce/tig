package domain

import "github.com/byjayce/tig/internal/file"

// Command is common interface for all commands
type Command[Param any] func(param Param) error

// TigRepository is API for tig repository.
type TigRepository interface {
	// IsInitialized returns true if repository is already initialized.
	IsInitialized() bool
	Init(config file.Config, head file.SymbolicReference) error
}

type File interface {
	Path() string
	Bytes() []byte
}
