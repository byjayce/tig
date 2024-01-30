package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("HashObject", func() {
	var (
		baseDir string
		data    []byte
	)

	BeforeEach(func() {
		t := GinkgoT()
		baseDir = t.TempDir()
		data = []byte("hello world")
	})

	When("HashObject가 호출된 경우", func() {
		Context("DryRun이 false인 경우", func() {
			It("키를 리턴하고 에러는 nil이어야 한다.", func() {
				k, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    data,
				})
				Expect(err).To(BeNil())
				Expect(k).To(Equal(newKey(BlobType, data).String()))
				_, err = os.Stat(newKey(BlobType, data).Path(baseDir))
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("DryRun이 true인 경우", func() {
			It("실제 생성하지 않고, 키를 리턴하고 에러는 nil이어야 한다.", func() {
				k, err := HashObject(HashObjectParam{
					DryRun:  true,
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    data,
				})
				expectedKey := newKey(BlobType, data)

				Expect(err).To(BeNil())
				Expect(k).To(Equal(expectedKey.String()))
				_, err = os.Stat(expectedKey.Path(baseDir))
				Expect(err).To(HaveOccurred())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})
})
