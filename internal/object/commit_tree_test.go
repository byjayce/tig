package object

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("CommitTree", func() {
	var (
		t           GinkgoTInterface
		workingCopy string
		tigDir      string
	)

	BeforeEach(func() {
		t = GinkgoT()
		workingCopy = t.TempDir()
		tigDir = filepath.Join(workingCopy, ".tig")
	})

	When("커밋할 트리 객체가 없는 경우", func() {
		It("에러가 발생한다.", func() {
			commit, err := CommitTree(CommitTreeParam{
				TigDir:     tigDir,
				ObjectHash: "not_exist_object_hash",
				Author:     "jayce",
				Message:    "test message",
			})

			Expect(err).Should(HaveOccurred())
			Expect(commit).To(BeEmpty())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})
	When("트리 객체가 아닌 경우", func() {
		var (
			objectHash string
		)

		BeforeEach(func() {
			hash, err := HashObject(HashObjectParam{
				TigDir: tigDir,
				Type:   Blob,
				Data:   []byte("hello world"),
			})
			if err != nil {
				t.Fatal(err)
			}

			objectHash = hash
		})

		It("에러가 발생한다.", func() {
			commit, err := CommitTree(CommitTreeParam{
				TigDir:     tigDir,
				ObjectHash: objectHash,
				Message:    "test message",
				Author:     "jayce",
			})
			Expect(err).Should(HaveOccurred())
			Expect(commit).To(BeEmpty())
			Expect(err.Error()).To(Equal(fmt.Sprintf("%s is not a valid 'tree' object", objectHash)))
		})
	})
	When("트리 객체를 잘 전달한 경우", func() {
		var (
			objectHash string
		)

		BeforeEach(func() {
			hash, err := HashObject(HashObjectParam{
				TigDir: tigDir,
				Type:   Tree,
				Data:   []byte("tree"),
			})
			if err != nil {
				t.Fatal(err)
			}

			objectHash = hash
		})

		It("커밋을 생성한다.", func() {
			commit, err := CommitTree(CommitTreeParam{
				TigDir:     tigDir,
				ObjectHash: objectHash,
				Message:    "test commit",
				Author:     "jayce",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(commit).NotTo(BeEmpty())
		})
	})
})
