package porcelain

import (
	"github.com/byjayce/tig/internal/config"
	"github.com/byjayce/tig/internal/object"
	"github.com/byjayce/tig/internal/reference"
	"path/filepath"
)

type TagParam struct {
	User    config.User
	TigDir  string
	TagName string
	Target  string // Target object hash | reference | symbolic reference
	Message string
	Delete  bool
}

func Tag(param TagParam) error {
	referencePath := filepath.Join("refs", "tags", param.TagName)
	if param.Delete {
		return reference.UpdateRef(reference.UpdateRefParam{
			TigDir:        param.TigDir,
			ReferencePath: referencePath,
			Delete:        true,
		})
	}

	if param.Target == "" {
		return nil
	}

	objHash, err := RevParse(RevParseParam{
		TigDir: param.TigDir,
		Target: param.Target,
	})
	if err != nil {
		return err
	}

	if param.Message == "" {
		// Lightweight 태그
		return reference.UpdateRef(reference.UpdateRefParam{
			TigDir:        param.TigDir,
			ReferencePath: referencePath,
			ObjectHash:    objHash,
		})
	}

	t, err := object.CatFile(object.CatFileParam{
		TigDir:        param.TigDir,
		OperationType: object.CatFileOperationTypeType,
		ObjectHash:    objHash,
	})

	if err != nil {
		return err
	}

	// Annotated 태그
	hash, err := object.MKTag(object.MKTagParam{
		TigDir:     param.TigDir,
		ObjectType: object.Type(t),
		ObjectHash: objHash,
		Name:       param.TagName,
		Message:    param.Message,
		Tagger:     param.User.Email,
	})
	if err != nil {
		return err
	}

	return reference.UpdateRef(reference.UpdateRefParam{
		TigDir:        param.TigDir,
		ReferencePath: referencePath,
		ObjectHash:    hash,
	})
}
