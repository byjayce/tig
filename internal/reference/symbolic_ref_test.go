package reference

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("SymbolicRef", func() {
	When("SymbolicRef를 호출할 때", func() {
		var (
			t       GinkgoTInterface
			baseDir string
		)

		BeforeEach(func() {
			t = GinkgoT()
			baseDir = t.TempDir()
		})

		Context("Delete 옵션이 false인 경우", func() {
			Context("Predefined 심볼릭 참조 타입이 아닌 경우", func() {
				It("심볼릭 참조를 업데이트한다.", func() {
					_, err := SymbolicRef(SymbolicRefParam{
						BaseDir:       baseDir,
						Type:          "custom",
						ReferencePath: "non",
					})

					Expect(err).To(BeNil())
					data, err := os.ReadFile(filepath.Join(baseDir, "custom"))
					Expect(err).To(BeNil())
					Expect(string(data)).To(Equal("ref: non"))
				})
			})

			Context("ReferencePath가 빈 문자열인 경우", func() {
				It("심볼릭 참조를 반환한다.", func() {
					if err := os.WriteFile(filepath.Join(baseDir, string(Head)), []byte("ref: refs/heads/main"), os.ModePerm); err != nil {
						t.Fatal(err)
					}

					data, err := SymbolicRef(SymbolicRefParam{
						BaseDir: baseDir,
						Type:    Head,
					})

					Expect(err).To(BeNil())
					Expect(data).To(Equal("refs/heads/main"))
				})
			})

			It("refs/로 시작하지 않는 참조 경로인 경우 에러를 반환한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					BaseDir:       baseDir,
					Type:          Head,
					ReferencePath: "non",
				})

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("non: refusing to point HEAD outside of refs/"))
			})

			It("심볼릭 참조를 업데이트한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					BaseDir:       baseDir,
					Type:          Head,
					ReferencePath: "refs/heads/main",
				})

				Expect(err).To(BeNil())
				data, err := os.ReadFile(filepath.Join(baseDir, string(Head)))
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal("ref: refs/heads/main"))
			})
		})

		Context("Delete 옵션이 true인 경우", func() {
			It("Predefined 심볼릭 참조 타입인 경우 에러를 반환한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					BaseDir: baseDir,
					Type:    Head,
					Delete:  true,
				})

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("HEAD: refusing to delete predefined symbolic ref"))
			})

			It("일반 심볼릭 참조를 삭제한다.", func() {
				const symbolicRefPath = "custom"

				if err := os.WriteFile(filepath.Join(baseDir, symbolicRefPath), []byte("ref: refs/heads/main"), os.ModePerm); err != nil {
					t.Fatal(err)
				}

				_, err := SymbolicRef(SymbolicRefParam{
					BaseDir: baseDir,
					Type:    symbolicRefPath,
					Delete:  true,
				})

				Expect(err).To(BeNil())
				_, err = os.Stat(filepath.Join(baseDir, symbolicRefPath))
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
