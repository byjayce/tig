package object

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type WriteTreeParam struct {
	TigDir string
}

func WriteTree(param WriteTreeParam) (string, error) {
	// 인덱스 파일을 찾기
	idx, err := openIndex(param.TigDir)
	if err != nil {
		return "", err
	}
	// 인덱스 파일의 엔트리를 기반으로 루트트리 (WorkingCopy의 구조를 트리로 만듦)
	tree, err := buildTree(idx)
	if err != nil {
		return "", err
	}

	// Tree 객체 만들기
	return createTreeObject(param.TigDir, tree)
}

func createTreeObject(tigDir string, tree map[string]any) (string, error) {
	var (
		workingCopy = filepath.Dir(tigDir)
	)

	entries := make(TreeEntries, 0, len(tree))
	for key, value := range tree {
		subTree, ok := value.(map[string]any)
		if ok {
			hash, err := createTreeObject(tigDir, subTree)
			if err != nil {
				return "", err
			}

			stat, err := os.Stat(filepath.Join(workingCopy, key))
			if err != nil {
				return "", err
			}

			entries = append(entries, TreeEntry{
				Mode:       stat.Mode(),
				File:       key,
				ObjectHash: hash,
			})
			continue
		}

		entry, ok := value.(TreeEntry)
		if !ok {
			return "", errors.New("invalid values in tree")
		}

		filePath := filepath.Join(workingCopy, entry.File)

		// 인덱스에 추가된 파일이 실제 존재하는가 확인, 파일 권한 확인
		stat, err := os.Stat(filePath)
		if err != nil {
			return "", err
		}

		// 파일 읽어서 Blob 객체 만들기
		fileData, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}

		hash, err := HashObject(HashObjectParam{
			TigDir: tigDir,
			Type:   Blob,
			Data:   fileData,
		})
		if err != nil {
			return "", err
		}

		entries = append(entries, TreeEntry{
			Mode:       stat.Mode(),
			File:       entry.File,
			ObjectHash: hash,
		})
	}

	sort.Sort(entries)
	data, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}

	return HashObject(HashObjectParam{
		TigDir: tigDir,
		Type:   Tree,
		Data:   data,
	})
}

func buildTree(idx index) (map[string]any, error) {
	var (
		tree = make(map[string]any)
	)

	for _, trackingFile := range idx {
		var (
			parent     = tree
			parentPath = ""
			paths      = strings.Split(trackingFile.File, string([]byte{os.PathSeparator}))
		)

		// Tree 만들기

		// dirPath1, dirPath2, fileName
		for _, dirPath := range paths[:len(paths)-1] {
			fullPath := filepath.Join(parentPath, dirPath)

			dir, ok := parent[fullPath]
			if ok {
				parent = dir.(map[string]any)
				parentPath = filepath.Join(parentPath, dirPath)
				continue
			}

			m := make(map[string]any)
			parent[fullPath] = m
			parent = m
			parentPath = fullPath
		}

		parent[trackingFile.File] = TreeEntry{
			File: trackingFile.File,
		}
	}

	return tree, nil
}
