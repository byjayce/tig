package workingcopy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWorkingCopy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkingCopy Suite")
}
