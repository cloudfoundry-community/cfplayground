package server_test

import (
	"net/http"

	"code.google.com/p/go.net/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/simonleung8/cfplayground/server"
	"github.com/simonleung8/cfplayground/server/fakes"
)

var _ = Describe("Server", func() {
	var (
		fake *fakes.FakeServerHandlers
	)

	fake = new(fakes.FakeServerHandlers)
	go Serve(fake)

	Context("Server base url '/'", func() {
		res, err := http.Get("http://localhost:8080")

		It("should return http status 200 at path '/' ", func() {
			Ω(err).ShouldNot(HaveOccurred())

			Ω(res.StatusCode).To(Equal(200))
		})

		It("should called RedirectBase() at '/' ", func() {
			Ω(err).ShouldNot(HaveOccurred())

			Ω(fake.RedirectBaseCallCount()).To(Equal(1))
		})
	})

	It("should called initSession() at '/ws' ", func() {
		_, err := http.Get("http://localhost:8080/ws")
		Ω(err).ShouldNot(HaveOccurred())

		Ω(fake.InitSessionCallCount()).To(Equal(1))
	})

	//test the handlers in handlers_test
	PIt("should accept websocket at path /ws", func() {
		origin := "http://127.0.0.1/"
		wsUrl := "ws://localhost:8080/ws"

		//disabled gorilla websocket's origin check (origin == r.Host)
		//in order to accept the ws connection
		_, err := websocket.Dial(wsUrl, "ws", origin)
		Ω(err).ShouldNot(HaveOccurred())
	})

	PIt("should path '/' should redirect to '/ui", func() {
		res, err := http.Get("http://localhost:8080")
		Ω(err).ShouldNot(HaveOccurred())

		Ω(res.Request.URL.String()).Should(Equal("http://localhost:8080/ui/"))
	})

	PIt("should place uploaded files in user dir", func() {

	})

	PIt("should warn about invalid cf commands'", func() {

	})
})
