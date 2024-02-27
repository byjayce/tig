package object

import (
	"encoding/json"
	"fmt"
	"sort"
)

type LSTreeOption struct {
	Recursive bool // -r 옵션에 해당함
	Tree      bool // -t 옵션에 해당함
}

type LSTreeParam struct {
	TigDir     string
	ObjectHash string
	Option     LSTreeOption
}

func LSTree(param LSTreeParam) (TreeEntries, error) {
	// 객체를 파싱해야함
	obj, err := parseObject(param.TigDir, param.ObjectHash)
	if err != nil {
		return nil, err
	}

	// Validation (트리 객체가 맞냐 확인)
	if obj.Type != Tree {
		return nil, fmt.Errorf("%s is invalid tree hash", param.ObjectHash)
	}

	// TreeEntries 만들어서 순회
	var entries TreeEntries
	if err := json.Unmarshal(obj.Data, &entries); err != nil {
		return nil, err
	}

	var ret TreeEntries
	for _, e := range entries {
		// 옵션에 따라 리턴해줄 엔트리를 필터링
		if !e.Mode.IsDir() {
			ret = append(ret, e)
			continue
		}

		if param.Option.Tree {
			ret = append(ret, e)
		}

		if !param.Option.Recursive {
			continue
		}

		subEntries, err := LSTree(LSTreeParam{
			TigDir:     param.TigDir,
			ObjectHash: e.ObjectHash,
			Option:     param.Option,
		})
		if err != nil {
			return nil, err
		}

		ret = append(ret, subEntries...)
		continue
	}

	sort.Sort(ret)
	return ret, nil
}
