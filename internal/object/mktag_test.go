package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MKTag", func() {
	When("MKTag를 호출하면", func() {
		var (
			t       GinkgoTInterface
			baseDir string
		)

		BeforeEach(func() {
			t = GinkgoT()
			baseDir = t.TempDir()
		})

		It("태그 객체를 생성한다.", func() {
			hash, err := MKTag(MKTagParam{
				BaseDir:    baseDir,
				ObjectHash: "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
				Name:       "v1.0.0",
				Message:    "Initial commit",
				Tagger:     "byjayce",
			})

			Expect(err).Should(BeNil())
			Expect(hash).ShouldNot(BeEmpty())

			str, err := CatFile(CatFileParam{
				BaseDir:       baseDir,
				OperationType: CatFileOperationTypeType,
				Hash:          hash,
			})
			Expect(err).Should(BeNil())
			Expect(str).Should(Equal("tag"))
		})
	})
})
