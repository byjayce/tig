package workingcopy_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWorkingcopy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Workingcopy Suite")
}
