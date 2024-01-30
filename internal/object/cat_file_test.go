package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CatFile", func() {
	var (
		objectKey key
		data      []byte
		baseDir   string
	)

	BeforeEach(func() {
		t := GinkgoT()
		baseDir = t.TempDir()
		data = []byte("hello world")
		objectKey = newKey(BlobType, data)

		if _, err := HashObject(HashObjectParam{
			BaseDir: baseDir,
			Type:    BlobType,
			Data:    data,
		}); err != nil {
			t.Fatal(err)
		}
	})

	When("CatFile이 pretty-print으로 호출된 경우", func() {
		It("returns pretty-printed content", func() {
			c, err := CatFile(CatFileParam{
				BaseDir:       baseDir,
				OperationType: CatFileOperationTypePrettyPrint,
				Hash:          objectKey.String(),
			})
			Expect(err).To(BeNil())
			Expect(c).To(Equal(string(data)))
		})
	})
})
