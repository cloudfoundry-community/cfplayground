package users_test

import (
	"net/http"

	"github.com/cloudfoundry-community/cfplayground/cf/fakes"
	. "github.com/cloudfoundry-community/cfplayground/users"
	"github.com/cloudfoundry-community/cfplayground/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type responseWriter struct{}

func (r responseWriter) Header() http.Header       { return nil }
func (r responseWriter) Write([]byte) (int, error) { return 0, nil }
func (r responseWriter) WriteHeader(int)           {}

var _ = Describe("Users", func() {
	fakeCF := new(fakes.FakeCLI)
	user := New(
		responseWriter{},
		&http.Request{},
		"../",
		GenerateToken(),
		fakeCF,
		&websocket.Pipe{},
	)

	It("contains a token", func() {
		Ω(user.Token).ToNot(BeNil())
	})

	It("finds and returns the user with provided Token", func() {
		userWrong := List("Wrong_Token")
		Ω(user).ToNot(Equal(userWrong))

		userFound := List(user.Token)
		Ω(user).To(Equal(userFound))

	})
})
