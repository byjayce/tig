package reference

import (
	"errors"
	"github.com/byjayce/tig/internal/object"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("UpdateRef", func() {
	var (
		t           GinkgoTInterface
		workingCopy string
		tigDir      string
		objectHash  string
	)

	BeforeEach(func() {
		t = GinkgoT()
		workingCopy = t.TempDir()
		tigDir = filepath.Join(workingCopy, ".tig")

		hash, err := object.HashObject(object.HashObjectParam{
			TigDir: tigDir,
			Type:   object.Commit,
			Data:   []byte("test commit object"),
		})

		if err != nil {
			t.Fatal(err)
		}

		objectHash = hash
	})

	When("Delete 옵션이 false", func() {
		var (
			refPath = filepath.Join("refs", "heads", "test")
		)

		It("참조를 업데이트한다.", func() {
			err := UpdateRef(UpdateRefParam{
				TigDir:        tigDir,
				ReferencePath: refPath,
				ObjectHash:    objectHash,
			})

			Expect(err).ShouldNot(HaveOccurred())

			_, err = os.Stat(filepath.Join(tigDir, refPath))
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("Delete 옵션이 true", func() {
		var (
			refPath = filepath.Join("refs", "heads", "test")
		)

		It("참조를 삭제한다.", func() {
			err := UpdateRef(UpdateRefParam{
				TigDir:        tigDir,
				ReferencePath: refPath,
				ObjectHash:    objectHash,
				Delete:        true,
			})

			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat(filepath.Join(tigDir, refPath))
			Expect(err).Should(HaveOccurred())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})
	})
})
