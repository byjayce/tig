package porcelain_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPorcelain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Porcelain Suite")
}
