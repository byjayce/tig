package porcelain

import (
	"encoding/json"
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
			config := config.Config{
				Core: config.Core{
					Bare: false,
				},
			}

			initParam := InitParam{
				WorkingCopyPath: tempDir,
				Config:          config,
			}

			configBuf, err := json.Marshal(config)
			if err != nil {
				Fail(err.Error())
			}

			Expect(Init(initParam)).Should(BeNil())

			f, err := os.Stat(filepath.Join(tempDir, baseDir, configFileName))
			Expect(err).Should(BeNil())
			Expect(f.IsDir()).Should(BeFalse())

			buf, err := os.ReadFile(filepath.Join(tempDir, baseDir, configFileName))
			Expect(err).Should(BeNil())
			Expect(buf).Should(Equal(configBuf))

			f, err = os.Stat(filepath.Join(tempDir, baseDir, headFileName))
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
