package object

import (
	"encoding/json"
	"fmt"
	"time"
)

type commitObjectBody struct {
	Tree    string // Tree: Object Hash
	Parent  []string
	Author  string
	Date    time.Time
	Message string
}

type CommitTreeParam struct {
	TigDir           string
	ObjectHash       string
	ParentCommitHash []string
	Message          string
	Author           string
}

func CommitTree(param CommitTreeParam) (string, error) {
	v, err := parseObject(param.TigDir, param.ObjectHash)
	if err != nil {
		return "", err
	}

	if v.Type != Tree {
		return "", fmt.Errorf("%s is not a valid 'tree' object", param.ObjectHash)
	}

	body := commitObjectBody{
		Tree:    param.ObjectHash,
		Parent:  param.ParentCommitHash,
		Author:  param.Author,
		Date:    time.Now(),
		Message: param.Message,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return HashObject(HashObjectParam{
		TigDir: param.TigDir,
		Type:   Commit,
		Data:   data,
	})
}
