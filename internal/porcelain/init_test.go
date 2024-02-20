package porcelain

import (
	"github.com/byjayce/tig/internal/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Init", func() {
	When("Init() 함수를 호출하면", func() {
		var tempDir string
		BeforeEach(func() {
			tempDir = GinkgoT().TempDir()
		})

		It("config 파일과 head 파일이 생성된다.", func() {
			cfg := config.Config{
				Core: config.Core{
					Bare: false,
				},
				User: config.User{
					Name:  "jayce",
					Email: "jayce@byjayce.cc",
				},
			}

			initParam := InitParam{
				WorkingCopyDir: tempDir,
				Config:         cfg,
			}
			Expect(Init(initParam)).Should(BeNil())

			createdCfg, err := config.ReadConfigFile(filepath.Join(tempDir, baseDir))
			Expect(err).Should(BeNil())
			Expect(createdCfg).Should(Equal(cfg))

			f, err := os.Stat(filepath.Join(tempDir, baseDir, headFileName))
			Expect(err).Should(BeNil())
			Expect(f.IsDir()).Should(BeFalse())

			d, err := os.Stat(filepath.Join(tempDir, baseDir, objectsDirName))
			Expect(err).Should(BeNil())
			Expect(d.IsDir()).Should(BeTrue())

			d, err = os.Stat(filepath.Join(tempDir, baseDir, refsDirName))
			Expect(err).Should(BeNil())
			Expect(d.IsDir()).Should(BeTrue())
		})
	})
})
