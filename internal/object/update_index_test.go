package object

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io/fs"
	"os"
	"path/filepath"
)

var _ = Describe("UpdateIndex", func() {
	When("UpdateIndex가 Caches와 함께 호출된 경우", func() {
		var (
			t           GinkgoTInterface
			workingCopy string
			baseDir     string
		)

		BeforeEach(func() {
			t = GinkgoT()
			workingCopy = t.TempDir()
			baseDir = filepath.Join(workingCopy, ".git")
		})

		Context("그리고 Add Option이 false인 경우", func() {
			It("인덱스 안에 객체가 있는 경우 인덱스 업데이트", func() {
				const (
					testFile = "exist.txt"
				)

				existHash, err := HashObject(HashObjectParam{
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
					Hash: existHash,
				})
				if err := writeIndex(baseDir, idx); err != nil {
					t.Fatal(err)
				}

				newHash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("new"),
				})
				if err != nil {
					t.Fatal(err)
				}

				err = UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Caches: []*indexEntry{
						{
							Name: testFile,
							Mode: 100644,
							Hash: newHash,
						},
					},
				})

				Expect(err).ToNot(HaveOccurred())

				idx, err = openIndex(baseDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(idx)).To(Equal(1))
				Expect(idx[0].Name).To(Equal(testFile))
				Expect(idx[0].Mode).To(Equal(fs.FileMode(100644)))
				Expect(idx[0].Hash).To(Equal(newHash))
			})

			It("인덱스 안에 객체가 없기 때문에 에러 발생", func() {
				testFileHash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("hello world"),
				})
				if err != nil {
					t.Fatal(err)
				}

				err = UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Caches: []*indexEntry{
						{
							Name: "test.txt",
							Mode: 100644,
							Hash: testFileHash,
						},
					},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("test.txt: cannot add to the index"))
			})
		})

		Context("그리고 Add Option이 true인 경우", func() {
			It("추가하려는 파일이 객체 디렉토리에 없기 때문에 에러 발생", func() {
				err := UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Caches: []*indexEntry{
						{
							Name: "test.txt",
							Mode: 100644,
							Hash: "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
						},
					},
					Add: true,
				})

				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})

			It("객체 DB에 있는 경우 인덱스에 추가", func() {
				testFileHash, err := HashObject(HashObjectParam{
					BaseDir: baseDir,
					Type:    BlobType,
					Data:    []byte("new test file"),
				})
				if err != nil {
					t.Fatal(err)
				}

				err = UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Caches: []*indexEntry{
						{
							Name: "test.txt",
							Mode: 100644,
							Hash: testFileHash,
						},
					},
					Add: true,
				})

				Expect(err).ToNot(HaveOccurred())

				idx, err := openIndex(baseDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(idx)).To(Equal(1))
				Expect(idx[0].Name).To(Equal("test.txt"))
				Expect(idx[0].Mode).To(Equal(fs.FileMode(100644)))
				Expect(idx[0].Hash).To(Equal(testFileHash))

				_, err = os.Stat(filepath.Join(baseDir, "objects", testFileHash[:2], testFileHash[2:]))
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
	When("UpdateIndex가 Files와 함께 호출된 경우", func() {
		var (
			t           GinkgoTInterface
			workingCopy string
			baseDir     string
		)

		BeforeEach(func() {
			t = GinkgoT()
			workingCopy = t.TempDir()
			baseDir = filepath.Join(workingCopy, ".git")
		})

		It("파일이 없는 경우", func() {
			err := UpdateIndex(UpdateIndexParam{
				BaseDir: baseDir,
				Files: []string{
					"not-exist.txt",
				},
			})

			Expect(err).To(HaveOccurred())
			Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
		})

		It("디렉토리인 경우", func() {
			var (
				dirName = "dir"
				dirPath = filepath.Join(workingCopy, dirName)
			)
			if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
				t.Fatal(err)
			}

			err := UpdateIndex(UpdateIndexParam{
				BaseDir: baseDir,
				Files: []string{
					dirName,
				},
			})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(dirName + ": is a directory"))
		})

		Context("Add Option이 false인 경우", func() {
			It("인덱스에 파일이 있는 경우 인덱스 업데이트", func() {
				var (
					testFileName = "exist.txt"
					testFilePath = filepath.Join(workingCopy, testFileName)
				)

				existHash, err := HashObject(HashObjectParam{
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
					Name: testFileName,
					Mode: 0100644,
					Hash: existHash,
				})
				if err := writeIndex(baseDir, idx); err != nil {
					t.Fatal(err)
				}

				if err := os.WriteFile(testFilePath, []byte("new"), 0644); err != nil {
					t.Fatal(err)
				}

				err = UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Files: []string{
						testFileName,
					},
				})

				Expect(err).ToNot(HaveOccurred())

				idx, err = openIndex(baseDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(idx)).To(Equal(1))
				Expect(idx[0].Name).To(Equal(testFileName))
				Expect(idx[0].Mode).To(Equal(fs.FileMode(0644)))
				Expect(idx[0].Hash).ToNot(Equal(existHash))
			})
			It("인덱스에 파일이 없는 경우 에러 발생", func() {
				var (
					testFileName = "not-exist.txt"
					testFilePath = filepath.Join(workingCopy, testFileName)
				)
				if err := os.WriteFile(testFilePath, []byte("test"), 0100644); err != nil {
					t.Fatal(err)
				}

				err := UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Files: []string{
						testFileName,
					},
				})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(testFileName + ": cannot add to the index"))
			})
		})

		Context("Add Option이 true인 경우", func() {
			It("인덱스에 새로운 객체 추가", func() {
				var (
					testFileName = "new.txt"
					testFilePath = filepath.Join(workingCopy, testFileName)
				)
				if err := os.WriteFile(testFilePath, []byte("new"), 0100644); err != nil {
					t.Fatal(err)
				}

				err := UpdateIndex(UpdateIndexParam{
					BaseDir: baseDir,
					Files: []string{
						testFileName,
					},
					Add: true,
				})

				Expect(err).ToNot(HaveOccurred())

				idx, err := openIndex(baseDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(idx)).To(Equal(1))
				Expect(idx[0].Name).To(Equal(testFileName))
				Expect(idx[0].Mode).To(Equal(fs.FileMode(0644)))
				Expect(idx[0].Hash).ToNot(BeEmpty())
			})
		})
	})
})
