package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("LSTree", func() {
	var (
		t           GinkgoTInterface
		workingCopy string
		tigDir      string
		treeHash    string
		testCases   = [][]string{
			{filepath.Join("dir1", "file1"), "file1Content"},
			{filepath.Join("file2"), "file2Content"},
		}
	)

	BeforeEach(func() {
		t = GinkgoT()
		workingCopy = t.TempDir()
		tigDir = filepath.Join(workingCopy, ".tig")

		// dir1/file1, file2
		setIndex(t, tigDir, testCases)
		hash, err := WriteTree(WriteTreeParam{TigDir: tigDir})
		if err != nil {
			t.Fatal(err)
		}

		treeHash = hash
	})

	When("재귀적 옵션이 켜진 경우", func() {
		It("하위 디렉토리의 파일들을 모두 가져온다", func() {
			entries, err := LSTree(LSTreeParam{
				TigDir:     tigDir,
				ObjectHash: treeHash,
				Option: LSTreeOption{
					Recursive: true,
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(entries)).To(Equal(2))

			Expect(entries[0].File).To(Equal(testCases[0][0]))
			Expect(entries[0].Mode.IsDir()).NotTo(BeTrue())

			Expect(entries[1].File).To(Equal(testCases[1][0]))
			Expect(entries[1].Mode.IsDir()).NotTo(BeTrue())
		})
	})

	When("트리 옵션이 켜진 경우", func() {
		It("트리 객체를 포함해 리스트를 만든다", func() {
			entries, err := LSTree(LSTreeParam{
				TigDir:     tigDir,
				ObjectHash: treeHash,
				Option: LSTreeOption{
					Tree: true,
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(entries)).To(Equal(2))
			Expect(entries[0].Mode.IsDir()).To(BeTrue())
			Expect(entries[0].File).To(Equal("dir1"))
			Expect(entries[1].Mode.IsDir()).NotTo(BeTrue())
			Expect(entries[1].File).To(Equal("file2"))
		})
	})

})
