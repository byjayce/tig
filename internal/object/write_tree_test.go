package object

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

			setIndex(t, tigDir, testCases)
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
