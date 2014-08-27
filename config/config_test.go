package config_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/cloudfoundry-community/cfplayground/config"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	It("returns error if provided file path is not found", func() {
		_, err := New("./fakePath.json")
		Ω(err).Should(HaveOccurred())
		Ω(err.Error()).Should(ContainSubstring("no such file found:"))
	})

	Context("it populates the struct{}", func() {
		var (
			b       []byte
			f       *os.File
			configs *Config
		)

		BeforeEach(func() {
			//set tmp json and write to tmp file
			b = []byte(`{
				"server": {
					"url": "testUrl.com",
					"login": "admin",
					"pass": "pass",
					"org": "org",
					"space": "space"
				}
			}`)

			f = setupConfigFile(b)
			configs, _ = New(f.Name())
		})

		AfterEach(func() {
			os.Remove(path.Join("./", f.Name()))
		})

		It("populates server url, login and password", func() {
			Ω(configs.Server.Url).ShouldNot(BeEmpty())
			Ω(configs.Server.Login).ShouldNot(BeEmpty())
			Ω(configs.Server.Pass).ShouldNot(BeEmpty())
		})

	})

	Context("returns error if any required field is missing", func() {
		var f *os.File

		AfterEach(func() {
			os.Remove(path.Join("./", f.Name()))
		})

		It("requires server url", func() {
			b := []byte(`{
				"server": {
					"login": "admin",
					"pass" : "pass"
				}
			}`)

			f = setupConfigFile(b)
			_, err := New(f.Name())

			Ω(err.Error()).Should(ContainSubstring("Error: server url is missing"))
		})

		It("requires server login", func() {
			b := []byte(`{
				"server": {
					"url": "testurl.com",
					"pass" : "pass",
					"org": "org",
					"space": "space"

				}
			}`)

			f = setupConfigFile(b)
			_, err := New(f.Name())

			Ω(err.Error()).Should(ContainSubstring("Error: server login is missing"))
		})

		It("requires server password", func() {
			b := []byte(`{
				"server": {
					"url": "testurl.com",
					"login": "admin",
					"org": "org",
					"space": "space"

				}
			}`)

			f = setupConfigFile(b)
			_, err := New(f.Name())

			Ω(err.Error()).Should(ContainSubstring("Error: server password is missing"))
		})

		It("requires server org", func() {
			b := []byte(`{
				"server": {
					"url": "testurl.com",
					"login": "admin",
					"pass": "pass",
					"space": "space"

				}
			}`)

			f = setupConfigFile(b)
			_, err := New(f.Name())

			Ω(err.Error()).Should(ContainSubstring("Error: server organization is missing"))
		})

		It("requires server space", func() {
			b := []byte(`{
				"server": {
					"url": "testurl.com",
					"login": "admin",
					"pass": "pass",
					"org": "org"
				}
			}`)

			f = setupConfigFile(b)
			_, err := New(f.Name())

			Ω(err.Error()).Should(ContainSubstring("Error: server space is missing"))
		})

	})

})

func setupConfigFile(b []byte) *os.File {
	f, _ := ioutil.TempFile(".", "testConfig")
	l, err := f.Write(b)
	Ω(err).ShouldNot(HaveOccurred())
	Ω(l).Should(Equal(len(b)))
	return f
}
