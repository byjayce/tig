package porcelain

import (
	"github.com/byjayce/tig/internal/object"
	"github.com/byjayce/tig/internal/reference"
	"path/filepath"
)

type TagParam struct {
	TagName string
	Target  string // Target object hash | reference | symbolic reference
	Message string
	Delete  bool
}

func (t *Tig) Tag(param TagParam) error {
	referencePath := filepath.Join("refs", "tags", param.TagName)

	if param.Delete {
		return reference.UpdateRef(reference.UpdateRefParam{
			BaseDir:       t.baseDir,
			ReferencePath: referencePath,
			Delete:        true,
		})
	}

	if param.Target == "" {
		return nil
	}

	objHash, err := resolveRef(t.baseDir, param.Target)
	if err != nil {
		return err
	}

	if param.Message == "" {
		// Lightweight tag
		return reference.UpdateRef(reference.UpdateRefParam{
			BaseDir:       t.baseDir,
			ReferencePath: referencePath,
			ObjectHash:    objHash,
		})
	}

	hash, err := object.MKTag(object.MKTagParam{
		BaseDir:    t.baseDir,
		ObjectHash: objHash,
		Name:       param.TagName,
		Message:    param.Message,
	})
	if err != nil {
		return err
	}

	return reference.UpdateRef(reference.UpdateRefParam{
		BaseDir:       t.baseDir,
		ReferencePath: referencePath,
		ObjectHash:    hash,
	})
}
