package porcelain

import (
	"errors"
	"github.com/byjayce/tig/internal/config"
	"github.com/byjayce/tig/internal/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("Tag", func() {
	var (
		t           GinkgoTInterface
		workingCopy string
		tigDir      string
		user        = config.User{
			Name:  "jayce",
			Email: "jayce@byjayce.cc",
		}
	)

	BeforeEach(func() {
		t = GinkgoT()
		workingCopy = t.TempDir()
		tigDir = filepath.Join(workingCopy, ".tig")
	})

	When("Delete 옵션이 true인 경우", func() {
		var existTagName = "v1.0.0"
		BeforeEach(func() {
			t.Log("INSIDE")
			// .tig/refs/tags 디렉터리 생성
			if err := os.MkdirAll(filepath.Join(tigDir, "refs", "tags"), os.ModePerm); err != nil {
				t.Fatal(err)
			}

			// 미리 파일 생성
			if err := os.WriteFile(filepath.Join(tigDir, "refs", "tags", existTagName), []byte("test"), os.ModePerm); err != nil {
				t.Fatal(err)
			}
		})

		It("참조를 삭제한다.", func() {
			Expect(Tag(TagParam{
				TigDir:  tigDir,
				User:    user,
				TagName: existTagName,
				Target:  "hash-text",
				Delete:  true,
			})).To(BeNil())

			// .tig/refs/tags/v1.0.0 파일이 삭제 됐는지 확인
			Expect(filepath.Join(tigDir, "refs", "tags", existTagName)).NotTo(BeAnExistingFile())
		})
	})
	When("Delete 옵션이 false인 경우", func() {
		var (
			objHash string
		)
		BeforeEach(func() {
			if err := os.MkdirAll(filepath.Join(tigDir, "refs", "tags"), os.ModePerm); err != nil {
				t.Fatal(err)
			}

			hash, err := object.HashObject(object.HashObjectParam{
				TigDir: tigDir,
				Type:   object.Tree,
				Data:   []byte("test"),
			})
			if err != nil {
				t.Fatal(err)
			}

			objHash = hash
		})

		It("참조와 객체를 생성한다.", func() {
			Expect(Tag(TagParam{
				TigDir:  tigDir,
				User:    user,
				TagName: "v1.0.0",
				Target:  objHash,
				Message: "test",
			})).To(BeNil())

			// .tig/refs/tags/v1.0.0 파일이 생성 됐는지 확인
			Expect(filepath.Join(tigDir, "refs", "tags", "v1.0.0")).To(BeAnExistingFile())

			tagHash, err := RevParse(RevParseParam{
				TigDir: tigDir,
				Target: "v1.0.0",
			})
			Expect(err).To(BeNil())

			// 태그 객체가 생성됐는지 확인
			_, err = object.CatFile(object.CatFileParam{
				TigDir:        tigDir,
				OperationType: object.CatFileOperationTypeExist,
				ObjectHash:    tagHash,
			})
			Expect(err).ShouldNot(HaveOccurred())
		})

		Context("Message 필드가 빈 문자열인 경우", func() {
			It("참조만 생성한다.", func() {
				Expect(Tag(TagParam{
					TigDir:  tigDir,
					User:    user,
					TagName: "v1.0.0",
					Target:  objHash,
				})).To(BeNil())

				// .tig/refs/tags/v1.0.0 파일이 생성 됐는지 확인
				Expect(filepath.Join(tigDir, "refs", "tags", "v1.0.0")).To(BeAnExistingFile())

				tagHash, err := RevParse(RevParseParam{
					TigDir: tigDir,
					Target: "v1.0.0",
				})
				Expect(err).To(BeNil())

				// 태그 객체가 생성되지 않았는지 확인
				_, err = object.CatFile(object.CatFileParam{
					TigDir:        tigDir,
					OperationType: object.CatFileOperationTypeExist,
					ObjectHash:    tagHash,
				})
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("객체를 찾지 못한 경우", func() {
			It("에러를 반환한다.", func() {
				err := Tag(TagParam{
					TigDir:  tigDir,
					User:    user,
					TagName: "v1.0.0",
					Target:  "hash-text",
					Message: "test",
				})
				Expect(err).Should(HaveOccurred())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})
	})

})
