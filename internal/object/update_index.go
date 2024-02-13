package object

import (
	"fmt"
	"os"
	"path/filepath"
)

type UpdateIndexParam struct {
	Caches  []*indexEntry
	Files   []string
	BaseDir string
	Add     bool
}

func UpdateIndex(param UpdateIndexParam) error {
	idx, err := openIndex(param.BaseDir)
	if err != nil {
		return err
	}

	if len(param.Caches) != 0 {
		// `--cacheinfo` 옵션으로 인덱스를 업데이트하는 경우
		for _, cache := range param.Caches {
			index, entry := searchIndex(idx, cache.Name)
			if entry != nil {
				idx[index] = cache
				continue
			}

			if !param.Add {
				return fmt.Errorf("%s: cannot add to the index", cache.Name)
			}

			idx = append(idx, cache)
		}

		if err := writeIndex(param.BaseDir, idx); err != nil {
			return err
		}

		return nil
	}

	for _, file := range param.Files {
		filePath := filepath.Join(filepath.Dir(param.BaseDir), file)

		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return fmt.Errorf("%s: is a directory", file)
		}

		f, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		hash, err := HashObject(HashObjectParam{
			BaseDir: param.BaseDir,
			Type:    BlobType,
			Data:    f,
		})
		if err != nil {
			return err
		}

		entry := &indexEntry{
			Name: file,
			Mode: stat.Mode(),
			Hash: hash,
		}

		index, existEntry := searchIndex(idx, file)
		if existEntry != nil {
			idx[index] = &indexEntry{
				Name: file,
				Mode: stat.Mode(),
				Hash: hash,
			}
			continue
		}

		if !param.Add {
			return fmt.Errorf("%s: cannot add to the index", file)
		}

		idx = append(idx, entry)
	}

	return writeIndex(param.BaseDir, idx)
}
