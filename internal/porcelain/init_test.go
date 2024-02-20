package porcelain

import (
	"github.com/byjayce/tig/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Init", func() {
	When("Init 함수를 호출하면", func() {
		var tempDir string
		BeforeEach(func() {
			// 임시 디렉토리를 만들었습니다.
			tempDir = GinkgoT().TempDir()
		})

		It("config, HEAD, objects, refs 파일이 생성된다.", func() {
			cfg := config.Config{
				User: config.User{
					Name:  "jayce",
					Email: "jayce@byjayce.cc",
				},
			}

			param := InitParam{
				WorkingCopyDir: tempDir,
				Config:         cfg,
			}

			Expect(Init(param)).Should(BeNil())

			// 설정 파일이 의도된 대로 만들어졌는지 확인
			createdCfg, err := config.ReadConfigFile(filepath.Join(tempDir, baseDir))
			Expect(err).To(BeNil())
			Expect(createdCfg).Should(Equal(cfg))

			// HEAD 파일이 생성됐는가? & 디렉토리가 아니어야 한다
			f, err := os.Stat(filepath.Join(tempDir, baseDir, headFileName))
			Expect(err).To(BeNil())
			Expect(f.IsDir()).NotTo(BeTrue())

			// objects 디렉토리가 생성됐는가?
			d, err := os.Stat(filepath.Join(tempDir, baseDir, objectsDirName))
			Expect(err).To(BeNil())
			Expect(d.IsDir()).To(BeTrue())

			// Refs 디렉토리가 생성 됐는가?
			d, err = os.Stat(filepath.Join(tempDir, baseDir, refsDirName))
			Expect(err).To(BeNil())
			Expect(d.IsDir()).To(BeTrue())
		})
	})
})
