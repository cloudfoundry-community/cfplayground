package copy_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCopy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Copy Suite")
}
