package object

import "os"

type TreeEntry struct {
	Mode       os.FileMode
	File       string
	ObjectHash string
}

type TreeEntries []TreeEntry

func (t TreeEntries) Len() int {
	return len(t)
}

func (t TreeEntries) Less(i, j int) bool {
	return t[i].File < t[j].File
}

func (t TreeEntries) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
