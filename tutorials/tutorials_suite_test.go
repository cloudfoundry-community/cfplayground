package tutorials_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTutorials(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tutorials Suite")
}
