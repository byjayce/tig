package object

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
)

var _ = Describe("LSTree", func() {
	When("LSTree를 호출할 때", func() {
		var (
			t            GinkgoTInterface
			baseDir      string
			workingCopy  string
			rootTreeHash string
		)

		const (
			testDir  = "a/b"
			testFile = "exist.txt"
		)

		BeforeEach(func() {
			/**
			WorkingCopy에 저장소 초기화 & 필요한 파일 생성
			Index 파일에 추가
			Tree 객체 생성
			*/
			t = GinkgoT()
			workingCopy = t.TempDir()
			baseDir = filepath.Join(workingCopy, ".git")

			if err := os.MkdirAll(filepath.Join(workingCopy, "a", "b"), os.ModePerm); err != nil {
				t.Fatal(err)
			}

			if err := os.WriteFile(filepath.Join(workingCopy, testDir, testFile), []byte("exist"), os.ModePerm); err != nil {
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
				Name: filepath.Join(testDir, testFile),
				Mode: 100644,
				Hash: hash,
			})

			if err := writeIndex(baseDir, idx); err != nil {
				t.Fatal(err)
			}

			hash, err = WriteTree(WriteTreeParam{
				BaseDir: baseDir,
			})
			if err != nil {
				t.Fatal(err)
			}

			rootTreeHash = hash
		})

		Context("재귀적 옵션이 켜져있고", func() {
			Context("트리 옵션이 켜진 경우", func() {
				It("트리와 하위 모든 객체를 반환한다.", func() {
					entries, err := LSTree(LSTreeParam{
						BaseDir: baseDir,
						Hash:    rootTreeHash,
						Option: LSTreeOption{
							Recursive: true,
							Tree:      true,
						},
					})
					Expect(err).To(BeNil())
					Expect(entries).To(HaveLen(3))
					Expect(entries[0].Name).To(Equal("a"))
					Expect(entries[1].Name).To(Equal("a/b"))
					Expect(entries[2].Name).To(Equal("a/b/exist.txt"))
				})
			})

			Context("트리 옵션이 꺼진 경우", func() {
				It("트리를 제외한 하위 모든 객체를 반환한다.", func() {
					entries, err := LSTree(LSTreeParam{
						BaseDir: baseDir,
						Hash:    rootTreeHash,
						Option: LSTreeOption{
							Recursive: true,
							Tree:      false,
						},
					})
					Expect(err).To(BeNil())
					Expect(entries).To(HaveLen(1))
					Expect(entries[0].Name).To(Equal("a/b/exist.txt"))
				})
			})
		})

		Context("재귀적 옵션이 꺼져있고", func() {
			Context("트리 옵션이 켜진 경우", func() {
				It("하위 탐색하지 않고 모든 객체를 반환한다.", func() {
					entries, err := LSTree(LSTreeParam{
						BaseDir: baseDir,
						Hash:    rootTreeHash,
						Option: LSTreeOption{
							Recursive: false,
							Tree:      true,
						},
					})
					Expect(err).To(BeNil())
					Expect(entries).To(HaveLen(1))
					Expect(entries[0].Name).To(Equal("a"))
				})
			})

			Context("트리 옵션이 없는 경우", func() {
				It("하위 탐색하지 않고 트리를 제외한 모든 객체를 반환한다.", func() {
					entries, err := LSTree(LSTreeParam{
						BaseDir: baseDir,
						Hash:    rootTreeHash,
						Option: LSTreeOption{
							Recursive: false,
							Tree:      false,
						},
					})
					Expect(err).To(BeNil())
					Expect(entries).To(HaveLen(0))
				})
			})
		})
	})
})
