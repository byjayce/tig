package object

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("WriteTree", func() {
	When("WriteTree를 호출할 때", func() {
		Context("하위 트리를 만들 필요가 없는 경우", func() {
			const (
				testFile = "exist.txt"
			)

			var (
				t       GinkgoTInterface
				baseDir string
			)

			BeforeEach(func() {
				t = GinkgoT()
				baseDir = t.TempDir()

				hash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("exist"),
				})
				if err != nil {
					t.Fatal(err)
				}

				idx, err := openIndex(baseDir)
				if err != nil {
					t.Fatal(err)
				}

				idx = append(idx, &indexEntry{
					Name: testFile,
					Mode: 100644,
					Hash: hash,
				})

				if err := writeIndex(baseDir, idx); err != nil {
					t.Fatal(err)
				}
			})

			It("루트 트리 객체는 Blob 엔트리만 가진 형태로 트리 객체를 반환한다.", func() {
				hash, err := WriteTree(WriteTreeParam{
					BaseDir: baseDir,
				})
				Expect(err).To(BeNil())
				Expect(hash).ToNot(BeEmpty())

				file, err := CatFile(CatFileParam{
					BaseDir:       baseDir,
					OperationType: CatFileOperationTypePrettyPrint,
					Hash:          hash,
				})
				Expect(err).To(BeNil())

				var tree treeEntries
				if err := json.Unmarshal([]byte(file), &tree); err != nil {
					t.Fatal(err)
				}

				Expect(tree).To(HaveLen(1))
				Expect(tree[0].Name).To(Equal(testFile))
			})
		})
		Context("하위 트리를 만들 필요가 있는 경우", func() {
			const (
				testDir  = "dir"
				testFile = "dir/exist.txt"
			)
			var (
				t           GinkgoTInterface
				workingCopy string
				baseDir     string
			)

			BeforeEach(func() {
				t = GinkgoT()
				workingCopy = t.TempDir()
				baseDir = filepath.Join(workingCopy, ".git")
				if err := os.Mkdir(filepath.Join(workingCopy, testDir), 0755); err != nil {
					t.Fatal(err)
				}

				hash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("exist"),
				})
				if err != nil {
					t.Fatal(err)
				}

				idx, err := openIndex(baseDir)
				if err != nil {
					t.Fatal(err)
				}

				idx = append(idx, &indexEntry{
					Name: testFile,
					Mode: 100644,
					Hash: hash,
				})

				if err := writeIndex(baseDir, idx); err != nil {
					t.Fatal(err)
				}
			})

			It("루트 트리 객체는 하위 트리를 가진 형태로 트리 객체를 반환한다.", func() {
				hash, err := WriteTree(WriteTreeParam{
					BaseDir: baseDir,
				})
				Expect(err).To(BeNil())
				Expect(hash).ToNot(BeEmpty())

				treeFile, err := CatFile(CatFileParam{
					BaseDir:       baseDir,
					OperationType: CatFileOperationTypePrettyPrint,
					Hash:          hash,
				})
				Expect(err).To(BeNil())

				var tree treeEntries
				if err := json.Unmarshal([]byte(treeFile), &tree); err != nil {
					t.Fatal(err)
				}

				Expect(tree).To(HaveLen(1))
				Expect(tree[0].Name).To(Equal("dir"))

				blobFile, err := CatFile(CatFileParam{
					BaseDir:       baseDir,
					OperationType: CatFileOperationTypePrettyPrint,
					Hash:          tree[0].Hash,
				})
				Expect(err).To(BeNil())

				if err := json.Unmarshal([]byte(blobFile), &tree); err != nil {
					t.Fatal(err)
				}

				Expect(tree).To(HaveLen(1))
				Expect(tree[0].Name).To(Equal("exist.txt"))
			})
		})
	})
})
