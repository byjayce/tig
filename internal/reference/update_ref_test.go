package reference

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("UpdateRef", func() {
	When("UpdateRef를 호출할 때", func() {
		var (
			t       GinkgoTInterface
			baseDir string
		)

		BeforeEach(func() {
			t = GinkgoT()
			baseDir = t.TempDir()
		})

		Context("Delete 옵션이 false인 경우", func() {
			It("참조를 업데이트한다.", func() {
				const (
					filePath = "refs/heads/main"
					hash     = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
				)

				err := UpdateRef(UpdateRefParam{
					BaseDir:       baseDir,
					ReferencePath: filePath,
					ObjectHash:    hash,
				})
				Expect(err).To(BeNil())

				data, err := os.ReadFile(filepath.Join(baseDir, filePath))
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(hash))
			})
		})
		Context("Delete 옵션이 true인 경우", func() {
			It("참조를 삭제한다.", func() {
				const (
					filePath = "refs/heads/main"
					hash     = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"
				)

				if err := os.MkdirAll(filepath.Join(baseDir, filepath.Dir(filePath)), os.ModePerm); err != nil {
					t.Fatal(err)
				}

				if err := os.WriteFile(filepath.Join(baseDir, filePath), []byte(hash), os.ModePerm); err != nil {
					t.Fatal(err)
				}

				err := UpdateRef(UpdateRefParam{
					BaseDir:       baseDir,
					ReferencePath: filePath,
					Delete:        true,
				})
				Expect(err).To(BeNil())

				_, err = os.Stat(filepath.Join(baseDir, filePath))
				Expect(err).ToNot(BeNil())
				Expect(os.IsNotExist(err)).To(BeTrue())
			})
		})
	})
})
