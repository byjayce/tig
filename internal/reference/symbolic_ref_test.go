package reference

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("SymbolicRef", func() {
	var (
		t           GinkgoTInterface
		workingCopy string
		tigDir      string
	)

	BeforeEach(func() {
		t = GinkgoT()
		workingCopy = t.TempDir()
		tigDir = filepath.Join(workingCopy, ".tig")
		if err := os.MkdirAll(tigDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	})

	When("Delete 옵션이 false인 경우", func() {
		Context("Predefined 심볼릭 참조 타입이 아닌 경우", func() {
			It("심볼릭 참조를 업데이트한다.", func() {
				refPath := "example"

				_, err := SymbolicRef(SymbolicRefParam{
					TigDir:        tigDir,
					Type:          "custom",
					ReferencePath: refPath,
				})

				Expect(err).ShouldNot(HaveOccurred())
				data, err := os.ReadFile(filepath.Join(tigDir, "custom"))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal(fmt.Sprintf("ref: %s", refPath)))
			})
		})

		Context("ReferencePath가 빈 문자열인 경우", func() {
			var (
				refPath = "refs/heads/main"
			)

			BeforeEach(func() {
				if err := os.WriteFile(filepath.Join(tigDir, string(Head)), []byte(fmt.Sprintf("ref: %s", refPath)), os.ModePerm); err != nil {
					t.Fatal(err)
				}
			})

			It("심볼릭 참조를 반환한다.", func() {
				data, err := SymbolicRef(SymbolicRefParam{
					TigDir: tigDir,
					Type:   Head,
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).To(Equal(refPath))
			})
		})

		Context("ReferencePath가 refs/로 시작하지 않는 경우", func() {
			It("에러를 반환한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					TigDir:        tigDir,
					Type:          Head,
					ReferencePath: "invalid",
				})
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("ReferencePath가 refs/로 시작하는 경우", func() {
			It("심볼릭 참조를 업데이트 한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					TigDir:        tigDir,
					Type:          Head,
					ReferencePath: "refs/heads/main",
				})

				Expect(err).ShouldNot(HaveOccurred())
				data, err := os.ReadFile(filepath.Join(tigDir, string(Head)))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(data)).To(Equal("ref: refs/heads/main"))
			})
		})

	})

	When("Delete 옵션이 true인 경우", func() {
		Context("Predefined 심볼릭 참조 타입인 경우", func() {
			It("에러를 반환한다.", func() {
				_, err := SymbolicRef(SymbolicRefParam{
					TigDir: tigDir,
					Type:   Head,
					Delete: true,
				})

				Expect(err).Should(HaveOccurred())
			})
		})

		It("일반 심볼릭 참조를 삭제한다.", func() {
			var symbolicFilePath SymbolicRefType = "custom"
			if err := os.WriteFile(filepath.Join(tigDir, string(symbolicFilePath)), []byte("ref: refs/heads/main"), os.ModePerm); err != nil {
				t.Fatal(err)
			}

			_, err := SymbolicRef(SymbolicRefParam{
				TigDir: tigDir,
				Type:   symbolicFilePath,
				Delete: true,
			})

			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(filepath.Join(tigDir, string(symbolicFilePath)))
			Expect(err).Should(HaveOccurred())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})
})
