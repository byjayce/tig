package workingcopy

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
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
			config := InitOption{
				Core: InitCoreConfig{
					Bare: false,
				},
			}

			initParam := InitParam{
				WorkingCopyPath: tempDir,
				Option:          config,
			}

			configBuf, err := yaml.Marshal(config)
			if err != nil {
				Fail(err.Error())
			}

			Expect(Init(initParam)).Should(BeNil())

			f, err := os.Stat(filepath.Join(tempDir, configFileName))
			Expect(err).Should(BeNil())
			Expect(f.IsDir()).Should(BeFalse())

			buf, err := os.ReadFile(filepath.Join(tempDir, configFileName))
			Expect(err).Should(BeNil())
			Expect(buf).Should(Equal(configBuf))

			f, err = os.Stat(filepath.Join(tempDir, headFileName))
			Expect(err).Should(BeNil())
			Expect(f.IsDir()).Should(BeFalse())

			d, err := os.Stat(filepath.Join(tempDir, objectsDirName))
			Expect(err).Should(BeNil())
			Expect(d.IsDir()).Should(BeTrue())

			d, err = os.Stat(filepath.Join(tempDir, refsDirName))
			Expect(err).Should(BeNil())
			Expect(d.IsDir()).Should(BeTrue())
		})
	})
})
