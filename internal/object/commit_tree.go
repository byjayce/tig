package object

import (
	"encoding/json"
	"fmt"
	"time"
)

type commitBody struct {
	Tree    string
	Parent  []string
	Author  string
	Date    time.Time
	Message string
}

type CommitTreeParam struct {
	BaseDir          string
	TreeHash         string
	ParentCommitHash []string
	Message          string
	Author           string
}

// CommitTree commit 객체를 만드는 함수
func CommitTree(param CommitTreeParam) (string, error) {
	v, err := parseObject(param.BaseDir, param.TreeHash)
	if err != nil {
		return "", err
	}

	if v.Type != TreeType {
		return "", fmt.Errorf("%s: is not a tree type", param.TreeHash)
	}

	body := commitBody{
		Tree:    param.TreeHash,
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
		BaseDir: param.BaseDir,
		Type:    CommitType,
		Data:    data,
	})
}
