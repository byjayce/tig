package object

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("HashObject", func() {
	When("HashObject 호출", func() {
		var (
			tigDir string
			data   = []byte("hello jayce")
		)

		BeforeEach(func() {
			t := GinkgoT()
			tigDir = filepath.Join(t.TempDir(), ".tig")
		})

		Context("DryRun = false", func() {
			It("키를 리턴, 에러는 nil, 실제로 객체 파일 만듦.", func() {
				k, err := HashObject(HashObjectParam{
					DryRun: false,
					TigDir: tigDir,
					Type:   Blob,
					Data:   data,
				})

				// 에러가 없는지 확인
				Expect(err).NotTo(HaveOccurred())

				// 의도한대로 키가 나왔는지 확인
				expectedKey := newKey(Blob, data)
				Expect(k).To(Equal(string(expectedKey)))

				// 객체 파일이 있는지 확인
				stat, err := os.Stat(Key(k).Path(tigDir))
				Expect(err).NotTo(HaveOccurred())
				Expect(stat.IsDir()).NotTo(BeTrue())
			})
		})

		Context("DryRun = true", func() {
			It("키를 리턴, 에러는 nil, 실제 객체 파일은 없음", func() {
				k, err := HashObject(HashObjectParam{
					DryRun: true,
					TigDir: tigDir,
					Type:   Blob,
					Data:   data,
				})

				// 에러가 없는지 확인
				Expect(err).NotTo(HaveOccurred())

				// 의도한대로 키가 나왔는지 확인
				expectedKey := newKey(Blob, data)
				Expect(k).To(Equal(string(expectedKey)))

				// 객체 파일이 없는지 확인
				_, err = os.Stat(Key(k).Path(tigDir))
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})
	})
})
