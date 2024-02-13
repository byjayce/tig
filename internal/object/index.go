package object

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type indexEntry struct {
	Mode os.FileMode
	Name string
	Hash string
}

type indexEntries []*indexEntry

func (e indexEntries) Len() int {
	return len(e)
}

func (e indexEntries) Less(i, j int) bool {
	return e[i].Name < e[j].Name
}

func (e indexEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

const indexFileName = "index"

func openIndex(baseDir string) (indexEntries, error) {
	filePath := filepath.Join(baseDir, indexFileName)
	buf, err := os.ReadFile(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		return indexEntries{}, nil
	}

	var entries indexEntries
	if err := json.Unmarshal(buf, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func writeIndex(baseDir string, entries indexEntries) error {
	sort.Sort(entries)
	buf, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	filePath := filepath.Join(baseDir, indexFileName)
	return os.WriteFile(filePath, buf, 0644)
}

func searchIndex(entries indexEntries, name string) (int, *indexEntry) {
	// binary search
	low, high := 0, len(entries)
	for low < high {
		mid := low + (high-low)/2
		if entries[mid].Name < name {
			low = mid + 1
		} else {
			high = mid
		}
	}

	if low < len(entries) && entries[low].Name == name {
		return low, entries[low]
	}

	return 0, nil
}
