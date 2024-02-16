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

var _ = Describe("(*Tig).Tag", func() {
	When("(*Tig).Tag를 호출할 때", func() {
		var (
			t   GinkgoTInterface
			tig *Tig
		)

		BeforeEach(func() {
			t = GinkgoT()
			baseDir := t.TempDir()
			if err := os.Chdir(baseDir); err != nil {
				t.Fatal(err)
			}

			if err := Init(InitParam{
				WorkingCopyPath: baseDir,
				Config: config.Config{
					Core: config.Core{
						Bare: false,
					},
					User: config.User{
						Name:  "jayce",
						Email: "jayce@byjayce.cc",
					},
				},
			}); err != nil {
				t.Fatal(err)
			}

			var err error
			tig, err = NewTig()
			if err != nil {
				t.Fatal(err)
			}
		})

		Context("Delete 옵션이 true인 경우", func() {
			It("참조를 삭제한다.", func() {
				Expect(tig.Tag(TagParam{
					TagName: "v1.0.0",
					Target:  "hash-text",
					Delete:  true,
				})).To(BeNil())

				if err := os.WriteFile(".git/refs/tags/v1.0.0", []byte("test"), os.ModePerm); err != nil {
					t.Fatal(err)
				}

				Expect(tig.Tag(TagParam{
					TagName: "v1.0.0",
					Target:  "hash-text",
					Delete:  true,
				})).To(BeNil())

				_, err := os.Stat(".git/refs/tags/v1.0.0")
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})

		Context("Target이 빈 문자열인 경우", func() {
			It("참조를 업데이트하지 않는다.", func() {
				Expect(tig.Tag(TagParam{
					TagName: "v1.0.0",
				})).To(BeNil())
			})
		})

		Context("객체를 찾지 못한 경우", func() {
			It("에러를 반환한다.", func() {
				err := tig.Tag(TagParam{
					TagName: "v1.0.0",
					Target:  "hash-text",
				})

				Expect(err).ToNot(BeNil())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})

		Context("Message가 빈 문자열인 경우", func() {
			var (
				hash string
			)
			BeforeEach(func() {
				// 타겟 객체 생성
				var err error
				hash, err = object.HashObject(object.HashObjectParam{
					BaseDir: tig.baseDir,
					Type:    object.BlobType,
					Data:    []byte("test"),
				})
				if err != nil {
					t.Fatal(err)
				}
			})

			It("Lightweight tag를 생성한다.", func() {
				Expect(tig.Tag(TagParam{
					TagName: "v1.0.0",
					Target:  hash,
				})).To(BeNil())

				data, err := os.ReadFile(".git/refs/tags/v1.0.0")
				Expect(err).To(BeNil())
				Expect(string(data)).To(Equal(hash))
			})
		})

		Context("그 외", func() {
			var (
				objectHash string
			)

			BeforeEach(func() {
				var err error
				objectHash, err = object.HashObject(object.HashObjectParam{
					BaseDir: tig.baseDir,
					Type:    object.BlobType,
					Data:    []byte("test"),
				})
				if err != nil {
					t.Fatal(err)
				}
			})

			It("태그 객체와 참조를 생성한다.", func() {
				err := tig.Tag(TagParam{
					TagName: "v1.0.0",
					Target:  objectHash,
					Message: "test",
				})
				Expect(err).To(BeNil())

				data, err := os.ReadFile(".git/refs/tags/v1.0.0")
				Expect(err).To(BeNil())
				Expect(string(data)).ToNot(BeEmpty())

				_, err = os.Stat(filepath.Join(tig.baseDir, "objects", string(data[:2]), string(data[2:])))
				Expect(err).To(BeNil())
			})
		})
	})
})
