package object

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("CommitTree", func() {
	When("CommitTree를 호출할 때", func() {
		var (
			t       GinkgoTInterface
			baseDir string
		)

		BeforeEach(func() {
			t = GinkgoT()
			baseDir = t.TempDir()
		})

		Context("커밋할 트리 객체가 없는 경우", func() {
			It("에러를 반환한다.", func() {
				commit, err := CommitTree(CommitTreeParam{
					BaseDir:  baseDir,
					TreeHash: "4b825dc642cb6eb9a060e54bf8d69288fbee4904",
					Author:   "byjayce",
					Message:  "Initial commit",
				})

				Expect(err).ShouldNot(BeNil())
				Expect(commit).Should(BeEmpty())
				Expect(errors.Is(err, os.ErrNotExist)).Should(BeTrue())
			})
		})
		Context("트리 객체가 아닌 경우", func() {
			var (
				objectHash string
			)
			BeforeEach(func() {
				hash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("exist"),
				})
				if err != nil {
					t.Fatal(err)
				}

				objectHash = hash
			})

			It("에러를 반환한다.", func() {
				commit, err := CommitTree(CommitTreeParam{
					BaseDir:  baseDir,
					TreeHash: objectHash,
					Author:   "byjayce",
					Message:  "Initial commit",
				})

				Expect(err).ShouldNot(BeNil())
				Expect(commit).Should(BeEmpty())
				Expect(err.Error()).Should(ContainSubstring("not a tree"))
			})
		})
		Context("트리 객체를 전달한 경우", func() {
			var (
				treeHash string
			)
			BeforeEach(func() {
				hash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    TreeType,
					Data:    []byte("exist"),
				})
				if err != nil {
					t.Fatal(err)
				}

				treeHash = hash
			})

			It("커밋을 생성한다.", func() {
				commit, err := CommitTree(CommitTreeParam{
					BaseDir:  baseDir,
					TreeHash: treeHash,
					Author:   "byjayce",
					Message:  "Initial commit",
				})

				Expect(err).Should(BeNil())
				Expect(commit).ShouldNot(BeEmpty())
			})
		})
	})
})
