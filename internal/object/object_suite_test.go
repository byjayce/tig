package object

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestObjects(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Object Suite")
}

func setIndex(t GinkgoTInterface, tigDir string, testCases [][]string) {
	var (
		workingCopy = filepath.Dir(tigDir)
	)

	idx, err := openIndex(tigDir)
	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range testCases {
		// 테스트 케이스의 파일을 실제로 만들어주고, Blob 객체 만들기
		if err := os.MkdirAll(filepath.Join(workingCopy, filepath.Dir(tc[0])), os.ModePerm); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(filepath.Join(workingCopy, tc[0]), []byte(tc[1]), os.ModePerm); err != nil {
			t.Fatal(err)
		}

		hash, err := HashObject(HashObjectParam{
			TigDir: tigDir,
			Type:   Blob,
			Data:   []byte(tc[1]),
		})

		if err != nil {
			t.Fatal(err)
		}

		testCases[i] = append(tc, hash)

		idx = append(idx, &indexEntry{TreeEntry{
			Mode:       os.ModePerm,
			File:       tc[0],
			ObjectHash: hash,
		}})
	}

	if err := writeIndex(tigDir, idx); err != nil {
		t.Fatal(err)
	}
}
