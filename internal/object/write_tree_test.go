package object

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("WriteTree", func() {
	When("WriteTree 호출", func() {
		var (
			t           GinkgoTInterface
			workingCopy string
			tigDir      string
		)

		var testCases = [][]string{
			{filepath.Join("dir1", "file1"), "file1Content"},
			{"file2", "file2Content"},
		}

		BeforeEach(func() {
			t = GinkgoT()
			workingCopy = t.TempDir()
			tigDir = filepath.Join(workingCopy, ".tig")

			// 인덱스 파일 미리 준비
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
		})

		It("루트 트리 아래, 서브 트리 및 Blob 객체를 만들고 루트 트리 객체 해시를 전달한다", func() {
			// dir1 -> 서브 디렉토리, dir1/file1, file2 -> Blob 객체
			treeHash, err := WriteTree(WriteTreeParam{TigDir: tigDir})
			Expect(err).NotTo(HaveOccurred())
			Expect(treeHash).NotTo(Equal(""))

			objData, err := CatFile(CatFileParam{
				TigDir:        tigDir,
				OperationType: CatFileOperationTypePrettyPrint,
				ObjectHash:    treeHash,
			})
			Expect(err).NotTo(HaveOccurred())

			var entries TreeEntries
			if err := json.Unmarshal([]byte(objData), &entries); err != nil {
				t.Fatal(err)
			}

			Expect(len(entries)).To(Equal(2))
			Expect(entries[0].File).To(Equal("dir1"))
			Expect(entries[1].File).To(Equal("file2"))
		})
	})
})
