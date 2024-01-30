package object

import "fmt"

func content(t Type, data []byte) string {
	return fmt.Sprintf("%s %d\000%s", t, len(data), data)
}
