package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CatFile", func() {
	var (
		t      GinkgoTInterface
		tigDir string
	)

	BeforeEach(func() {
		t = GinkgoT()
		tigDir = t.TempDir()
	})

	When("CatFile이 pretty-print로 호출된 경우", func() {
		It("객체의 콘텐츠가 출력됨", func() {
			data := []byte("hello jayce")
			hash, err := HashObject(HashObjectParam{
				DryRun: false,
				TigDir: tigDir,
				Type:   Blob,
				Data:   data,
			})
			if err != nil {
				t.Fatal(err)
			}

			c, err := CatFile(CatFileParam{
				TigDir:        tigDir,
				OperationType: CatFileOperationTypePrettyPrint,
				ObjectHash:    hash,
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(c).To(Equal(string(data)))
		})
	})
})
