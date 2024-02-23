package object

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("UpdateIndex", func() {
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

	When("Caches 필드를 사용하는 경우", func() {
		Context("Add = false", func() {
			Context("Index 파일 안에 객체가 있는 경우", func() {
				var (
					fileName = "exist-file"
				)
				BeforeEach(func() {
					idx, err := openIndex(tigDir)
					if err != nil {
						t.Fatal(err)
					}

					hash, err := HashObject(HashObjectParam{
						TigDir: tigDir,
						Type:   Blob,
						Data:   []byte("hello world"),
					})
					if err != nil {
						t.Fatal(err)
					}

					idx = append(idx, &indexEntry{TreeEntry{
						Mode:       os.ModePerm,
						File:       fileName,
						ObjectHash: hash,
					}})
					if err := writeIndex(tigDir, idx); err != nil {
						t.Fatal(err)
					}
				})

				It("새로운 객체가 인덱스 파일에 업데이트 된다.", func() {
					const changedHash = "random-sha1"
					err := UpdateIndex(UpdateIndexParam{
						Caches: []*TreeEntry{
							{
								Mode:       os.ModePerm,
								File:       fileName,
								ObjectHash: changedHash,
							},
						},
						TigDir: tigDir,
					})

					Expect(err).NotTo(HaveOccurred())

					idx, err := openIndex(tigDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(idx)).To(Equal(1))

					entryIdx, entry := searchIndex(idx, fileName)
					Expect(entry).NotTo(BeNil())
					Expect(entryIdx).To(Equal(0))
					Expect(entry.ObjectHash).To(Equal(changedHash))
				})
			})

			Context("Index 파일 안에 객체가 없는 경우", func() {
				var (
					fileName   = "not-exist.txt"
					objectHash = "not_exist_hash_value"
				)

				It("Index 업데이트를 할 수 없다는 에러를 반환한다.", func() {
					param := UpdateIndexParam{
						Caches: []*TreeEntry{
							{
								Mode:       os.ModePerm,
								File:       fileName,
								ObjectHash: objectHash,
							},
						},
						TigDir: tigDir,
					}

					err := UpdateIndex(param)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("%s: cannot add to the index", fileName)))
				})
			})
		})

		Context("Add = true", func() {
			Context("Index 파일 안에 객체가 없는 경우", func() {
				var (
					fileName   = "not-exist.txt"
					objectHash = "not_exist_hash_value"
				)

				It("Index 파일에 새로운 객체를 추가한다.", func() {
					err := UpdateIndex(UpdateIndexParam{
						Files: []string{},
						Caches: []*TreeEntry{
							{
								Mode:       os.ModePerm,
								File:       fileName,
								ObjectHash: objectHash,
							},
						},
						TigDir: tigDir,
						Add:    true,
					})

					Expect(err).NotTo(HaveOccurred())

					idx, err := openIndex(tigDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(idx)).To(Equal(1))
					Expect(idx[0].ObjectHash).To(Equal(objectHash))
					Expect(idx[0].File).To(Equal(fileName))
				})
			})
		})
	})

	When("Files 필드를 사용하는 경우", func() {
		Context("대상 파일이 없는 경우", func() {
			var (
				fileName = "not-exist.txt"
			)

			It("에러를 반환한다", func() {
				param := UpdateIndexParam{
					Files:  []string{fileName},
					TigDir: tigDir,
				}

				err := UpdateIndex(param)
				Expect(err).To(HaveOccurred())
				Expect(errors.Is(err, os.ErrNotExist)).To(BeTrue())
			})
		})

		Context("Add = false", func() {
			Context("Index 파일 안에 객체가 있는 경우", func() {
				var (
					fileName     = "exist-file"
					originalHash string
				)
				BeforeEach(func() {
					idx, err := openIndex(tigDir)
					if err != nil {
						t.Fatal(err)
					}

					hash, err := HashObject(HashObjectParam{
						TigDir: tigDir,
						Type:   Blob,
						Data:   []byte("hello world"),
					})
					if err != nil {
						t.Fatal(err)
					}

					originalHash = hash
					idx = append(idx, &indexEntry{TreeEntry{
						Mode:       os.ModePerm,
						File:       fileName,
						ObjectHash: hash,
					}})
					if err := writeIndex(tigDir, idx); err != nil {
						t.Fatal(err)
					}
				})

				It("새로운 객체가 인덱스 파일에 업데이트 된다.", func() {
					var (
						filePath = filepath.Join(workingCopy, fileName)
					)

					if err := os.WriteFile(filePath, []byte("hello jayce"), os.ModePerm); err != nil {
						t.Fatal(err)
					}

					err := UpdateIndex(UpdateIndexParam{
						Files: []string{
							fileName,
						},
						TigDir: tigDir,
					})
					Expect(err).NotTo(HaveOccurred())

					idx, err := openIndex(tigDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(idx)).To(Equal(1))

					entryIdx, entry := searchIndex(idx, fileName)
					Expect(entry).NotTo(BeNil())
					Expect(entryIdx).To(Equal(0))
					Expect(entry.ObjectHash).NotTo(Equal(originalHash))
				})
			})

			Context("Index 파일 안에 객체가 없는 경우", func() {
				var (
					fileName = "not-exist"
				)

				BeforeEach(func() {
					if err := os.WriteFile(filepath.Join(workingCopy, fileName), []byte("hello world"), os.ModePerm); err != nil {
						t.Fatal(err)
					}
				})

				It("Index 업데이트를 할 수 없다는 에러를 반환한다.", func() {
					param := UpdateIndexParam{
						Files:  []string{fileName},
						TigDir: tigDir,
					}

					err := UpdateIndex(param)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("%s: cannot add to the index", fileName)))
				})
			})
		})

		Context("Add = true", func() {
			Context("Index 파일 안에 객체가 없는 경우", func() {
				var (
					fileName = "not-exist"
				)

				BeforeEach(func() {
					if err := os.WriteFile(filepath.Join(workingCopy, fileName), []byte("hello world"), os.ModePerm); err != nil {
						t.Fatal(err)
					}
				})

				It("Index 파일에 새로운 객체를 추가한다.", func() {
					param := UpdateIndexParam{
						Add:    true,
						Files:  []string{fileName},
						TigDir: tigDir,
					}

					err := UpdateIndex(param)
					Expect(err).NotTo(HaveOccurred())

					idx, err := openIndex(tigDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(idx)).To(Equal(1))
					Expect(idx[0].File).To(Equal(fileName))
				})
			})
		})
	})
})
