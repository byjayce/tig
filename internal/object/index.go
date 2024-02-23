package object

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

const indexFileName = "index"

type indexEntry struct {
	TreeEntry
}

type index []*indexEntry

func (idx index) Len() int {
	return len(idx)
}

func (idx index) Less(i, j int) bool {
	return idx[i].File < idx[j].File
}

func (idx index) Swap(i, j int) {
	idx[i], idx[j] = idx[j], idx[i]
}

func openIndex(tigDir string) (index, error) {
	indexPath := filepath.Join(tigDir, indexFileName)
	data, err := os.ReadFile(indexPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return index{}, nil
	}

	var idx index
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}

	return idx, nil
}

func writeIndex(tigDir string, idx index) error {
	sort.Sort(idx)
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}

	indexPath := filepath.Join(tigDir, indexFileName)
	return os.WriteFile(indexPath, data, os.ModePerm)
}

func searchIndex(idx index, file string) (int, *indexEntry) {
	low, high := 0, len(idx)
	for low < high {
		mid := (low + high) / 2

		if idx[mid].File < file {
			low = mid + 1
			continue
		}

		high = mid
	}

	if low < len(idx) && idx[low].File == file {
		return low, idx[low]
	}

	return 0, nil
}
