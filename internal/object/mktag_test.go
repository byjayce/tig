package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"path/filepath"
)

var _ = Describe("MKTag", func() {
	When("MKTag를 호출하면", func() {
		var (
			t      GinkgoTInterface
			tigDir string
		)

		BeforeEach(func() {
			t = GinkgoT()
			workingCopy := t.TempDir()
			tigDir = filepath.Join(workingCopy, ".tig")
		})

		It("태그 객체를 생성한다.", func() {
			hash, err := MKTag(MKTagParam{
				TigDir:     tigDir,
				ObjectHash: "example-hash",
				Name:       "v1.0.0",
				Message:    "Initial commit",
				Tagger:     "byjayce",
			})

			Expect(err).Should(BeNil())
			Expect(hash).ShouldNot(BeEmpty())

			str, err := CatFile(CatFileParam{
				TigDir:        tigDir,
				OperationType: CatFileOperationTypeType,
				ObjectHash:    hash,
			})
			Expect(err).Should(BeNil())
			Expect(str).Should(Equal("tag"))
		})
	})
})
