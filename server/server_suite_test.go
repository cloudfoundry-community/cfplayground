package server_test

import (
	"os"
	"path"


	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	userToken string
)

func TestServer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
}

var _ = BeforeSuite(func() {
	userToken = "test_token"
	os.MkdirAll(path.Join("./", userToken, "app"), os.ModePerm)
})

var _ = AfterSuite(func() {
	err := os.RemoveAll(path.Join("./", userToken))
	Ω(err).ShouldNot(HaveOccurred())
})
