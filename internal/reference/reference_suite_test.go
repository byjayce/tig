package reference_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestReference(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reference Suite")
}
