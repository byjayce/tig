package object

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	TreeType Type = "tree"
)

type treeNode struct {
	Mode    os.FileMode
	Type    Type
	Hash    string
	Name    string
	entries map[string]*treeNode
}

type treeEntries []*treeNode

func (t treeEntries) Len() int {
	return len(t)
}

func (t treeEntries) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

func (t treeEntries) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func buildTreeEntry(idxEntries indexEntries) *treeNode {
	var (
		tree treeNode
	)

	tree.entries = make(map[string]*treeNode)

	for _, idxEntry := range idxEntries {
		segments := strings.Split(idxEntry.Name, string(os.PathSeparator))
		if len(segments) == 1 {
			tree.entries[segments[0]] = &treeNode{
				Mode: idxEntry.Mode,
				Type: BlobType,
				Hash: idxEntry.Hash,
				Name: segments[0],
			}
			continue
		}

		curr := &tree
		for _, segment := range segments[:len(segments)-1] {
			if _, ok := curr.entries[segment]; !ok {
				curr.entries[segment] = &treeNode{
					Type:    TreeType,
					Mode:    040000,
					Name:    segment,
					entries: make(map[string]*treeNode),
				}
			}

			curr = curr.entries[segment]
		}

		curr.entries[segments[len(segments)-1]] = &treeNode{
			Mode: idxEntry.Mode,
			Type: BlobType,
			Hash: idxEntry.Hash,
			Name: segments[len(segments)-1],
		}
	}

	return &tree
}

func writeTree(baseDir, treeDir string, entry *treeNode) (string, error) {
	if entry.Type == BlobType {
		return entry.Hash, nil
	}

	stat, err := os.Stat(treeDir)
	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", fmt.Errorf("%s: is not a directory", treeDir)
	}

	for name, node := range entry.entries {
		hash, err := writeTree(baseDir, filepath.Join(treeDir, name), node)
		if err != nil {
			return "", err
		}

		entry.entries[name] = &treeNode{
			Mode:    node.Mode,
			Type:    node.Type,
			Hash:    hash,
			Name:    name,
			entries: node.entries,
		}
	}

	var entries treeEntries
	for _, node := range entry.entries {
		entries = append(entries, node)
	}

	sort.Sort(entries)
	treeData, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}

	return HashObject(HashObjectParam{
		BaseDir: baseDir,
		Type:    TreeType,
		Data:    treeData,
	})
}

type WriteTreeParam struct {
	BaseDir string
}

func WriteTree(param WriteTreeParam) (string, error) {
	idx, err := openIndex(param.BaseDir)
	if err != nil {
		return "", err
	}

	tree := buildTreeEntry(idx)
	rootDir := filepath.Dir(param.BaseDir)
	return writeTree(param.BaseDir, rootDir, tree)
}
