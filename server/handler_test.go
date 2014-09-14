package server_test

import (
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	. "github.com/cloudfoundry-community/cfplayground/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Handler", func() {
	var (
		port     string
		listener net.Listener
	)
	//go Serve(NewHandler("../"))

	BeforeEach(func() {
		http.DefaultServeMux = http.NewServeMux()
		RegisterHandler(NewHandler("../"))
		listener, _ = net.Listen("tcp", "127.0.0.1:0")
		go func() {
			http.Serve(listener, nil)
		}()
		port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
	})

	AfterEach(func() {
		listener.Close()
	})

	Context("Deleting uploaded files", func() {
		var (
			dir string
			err error
		)

		BeforeEach(func() {
			dir, err = ioutil.TempDir("./", "")
			Ω(err).ShouldNot(HaveOccurred())
			err = os.MkdirAll(filepath.Join("./", dir, "containers", userToken, "app"), 0700)
			Ω(err).ShouldNot(HaveOccurred())
			err = os.Chdir(filepath.Join("./", dir))
			Ω(err).ShouldNot(HaveOccurred())
			_, err = ioutil.TempFile(filepath.Join("./", "containers", userToken, "app"), "tmp")
			Ω(err).ShouldNot(HaveOccurred())
			_, err = ioutil.TempFile(filepath.Join("./", "containers", userToken, "app"), "tmp")
			Ω(err).ShouldNot(HaveOccurred())
			_, err = ioutil.TempFile(filepath.Join("./", "containers", userToken, "app"), "tmp")
			Ω(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			os.Chdir("../")
			os.RemoveAll(filepath.Join("./", dir))
		})

		It("should delete all uploaded files with /delete endpoint", func() {
			res, err := http.Get("http://localhost:" + port + "/delete/" + userToken)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(res.StatusCode).To(Equal(200))

			fileInfo, err := ioutil.ReadDir(filepath.Join("./", "containers", userToken, "app"))
			Ω(err).To(HaveOccurred())
			Ω(len(fileInfo)).Should(Equal(0))
		})

	})
})
