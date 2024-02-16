package porcelain

import (
	"errors"
	"github.com/byjayce/tig/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("NewTag", func() {
	var (
		t       GinkgoTInterface
		tempDir string
	)

	BeforeEach(func() {
		t = GinkgoT()
		tempDir = t.TempDir()
	})

	When("NewTag를 호출할 때", func() {
		Context("설정 파일이 디렉토리 최상단에 있는 경우", func() {
			Context("그리고 bare 옵션이 꺼져있는 경우", func() {
				It("에러를 반환한다.", func() {
					// 작업 공간 이동
					if err := os.Chdir(tempDir); err != nil {
						t.Fatal(err)
					}

					// Bare 옵션을 끄고 최상단 디렉토리에 설정 파일을 생성한다.
					if err := config.CreateConfigFile(tempDir, config.Config{
						Core: config.Core{
							Bare: false,
						},
					}); err != nil {
						t.Fatal(err)
					}

					tig, err := NewTig()
					Expect(tig).To(BeNil())
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("not a git repository"))
				})
			})
			Context("그리고 bare 옵션이 켜져있는 경우", func() {
				It("Tig 구조체를 초기화한다.", func() {
					// 작업 공간 이동
					if err := os.Chdir(tempDir); err != nil {
						t.Fatal(err)
					}

					// Bare 옵션을 켜고 최상단 디렉토리에 설정 파일을 생성한다.
					if err := config.CreateConfigFile(tempDir, config.Config{
						Core: config.Core{
							Bare: true,
						},
					}); err != nil {
						t.Fatal(err)
					}

					tig, err := NewTig()
					Expect(tig).ToNot(BeNil())
					Expect(err).To(BeNil())
				})
			})
		})

		Context("설정 파일이 .git 디렉토리 아래 있는 경우", func() {
			It("Tig 구조체를 초기화한다.", func() {
				// 작업 공간 이동
				if err := os.Chdir(tempDir); err != nil {
					t.Fatal(err)
				}

				// .git 디렉토리를 생성하고 설정 파일을 생성한다.
				gitDir := tempDir + "/.git"
				if err := os.Mkdir(gitDir, os.ModePerm); err != nil {
					t.Fatal(err)
				}
				if err := config.CreateConfigFile(gitDir, config.Config{
					Core: config.Core{
						Bare: false,
					},
				}); err != nil {
					t.Fatal(err)
				}

				tig, err := NewTig()
				Expect(tig).ToNot(BeNil())
				Expect(err).To(BeNil())
			})
		})

		Context("설정 파일이 없는 경우", func() {
			It("에러를 반환한다.", func() {
				// 작업 공간 이동
				if err := os.Chdir(tempDir); err != nil {
					t.Fatal(err)
				}

				tig, err := NewTig()
				Expect(tig).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})
	})
})
