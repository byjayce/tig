package object

import (
	"fmt"
	"os"
	"path/filepath"
)

type TreeEntry struct {
	Mode       os.FileMode
	File       string
	ObjectHash string
}

type UpdateIndexParam struct {
	Files  []string
	Caches []*TreeEntry
	TigDir string
	Add    bool
}

func UpdateIndex(param UpdateIndexParam) error {
	// 인덱스 파일을 열고
	idx, err := openIndex(param.TigDir)
	if err != nil {
		return err
	}

	// 파라미터에 Caches로 동작하는지, Files로 동작하는지 분기
	if len(param.Caches) != 0 {
		// Caches로 동작하는 경우

		// Caches 순회하며 인덱스에서 검색
		for _, cache := range param.Caches {
			entryIdx, entry := searchIndex(idx, cache.File)
			if entry != nil {
				// 있는 경우는 새로운 캐시로 업데이트 해주기
				idx[entryIdx] = &indexEntry{
					TreeEntry: *cache,
				}
				continue
			}

			if !param.Add {
				// Add 옵션이 꺼져있는 경우 에러 반환
				return fmt.Errorf("%s: cannot add to the index", cache.File)
			}

			idx = append(idx, &indexEntry{TreeEntry: *cache})
		}

		if err := writeIndex(param.TigDir, idx); err != nil {
			return err
		}

		return nil
	}

	// 파일로 동작해야 하는 경우

	// 파일 경로를 순회하면서 실제 존재하는 파일인지 확인
	for _, file := range param.Files {
		// 파일 내용을 읽어서 Blob Object를 생성
		// tigDir == workingCopy/.tig
		filePath := filepath.Join(filepath.Dir(param.TigDir), file)
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.IsDir() {
			return fmt.Errorf("%s: is a directory", file)
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		hash, err := HashObject(HashObjectParam{
			TigDir: param.TigDir,
			Type:   Blob,
			Data:   data,
		})
		if err != nil {
			return err
		}

		// 생성된 해시를 기준으로 인덱스 서치
		entryIdx, entry := searchIndex(idx, file)
		if entry != nil {
			// 만약 인덱스에 있다면 업데이트 반영
			idx[entryIdx] = &indexEntry{TreeEntry{
				Mode:       stat.Mode(),
				File:       file,
				ObjectHash: hash,
			}}
			continue
		}

		if !param.Add {
			// 만약 Add 옵션이 꺼져있다면, 에러를 반환
			return fmt.Errorf("%s: cannot add to the index", file)
		}

		// 새롭게 해시를 추가
		idx = append(idx, &indexEntry{TreeEntry{
			Mode:       stat.Mode(),
			File:       file,
			ObjectHash: hash,
		}})
	}

	return writeIndex(param.TigDir, idx)
}
