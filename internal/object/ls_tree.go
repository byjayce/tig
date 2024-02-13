package object

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type entry struct {
	Mode os.FileMode
	Type Type
	Hash string
	Name string
}

type Entries []*entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Less(i, j int) bool {
	return e[i].Name < e[j].Name
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

type LSTreeOption struct {
	Recursive bool
	Tree      bool
}

type LSTreeParam struct {
	BaseDir string
	Hash    string
	Option  LSTreeOption
}

func LSTree(param LSTreeParam) (Entries, error) {
	v, err := parseObject(param.BaseDir, param.Hash)
	if err != nil {
		return nil, err
	}

	if v.Type != TreeType {
		return nil, fmt.Errorf("%s: not a tree object", param.Hash)
	}

	var treeData Entries
	if err := json.Unmarshal(v.Data, &treeData); err != nil {
		return nil, err
	}

	var ret Entries
	for _, e := range treeData {
		if e.Type == BlobType {
			ret = append(ret, e)
			continue
		}

		if param.Option.Tree {
			ret = append(ret, e)
		}

		if param.Option.Recursive {
			innerTree, err := LSTree(LSTreeParam{
				BaseDir: param.BaseDir,
				Hash:    e.Hash,
				Option:  param.Option,
			})
			if err != nil {
				return nil, err
			}

			for _, it := range innerTree {
				ret = append(ret, &entry{
					Mode: it.Mode,
					Type: it.Type,
					Hash: it.Hash,
					Name: filepath.Join(e.Name, it.Name),
				})
			}
		}
	}

	sort.Sort(ret)
	return ret, nil
}
