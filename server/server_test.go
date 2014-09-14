package server_test

import (
	"net"
	"net/http"
	"strconv"

	"code.google.com/p/go.net/websocket"
	. "github.com/cloudfoundry-community/cfplayground/server"
	"github.com/cloudfoundry-community/cfplayground/server/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		fake     *fakes.FakeServerHandlers
		port     string
		listener net.Listener
	)

	Context("Server base url '/'", func() {

		BeforeEach(func() {
			fake = new(fakes.FakeServerHandlers)
			http.DefaultServeMux = http.NewServeMux()
			RegisterHandler(fake)
			listener, _ = net.Listen("tcp", "127.0.0.1:0")
			go func() {
				http.Serve(listener, nil)
			}()
			port = strconv.Itoa(listener.Addr().(*net.TCPAddr).Port)
		})

		AfterEach(func() {
			listener.Close()
		})

		It("should called RedirectBase() at '/' ", func() {
			res, err := http.Get("http://localhost:" + port)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(res.StatusCode).To(Equal(200))
			Ω(fake.RedirectBaseCallCount()).To(Equal(1))
		})
	})

	It("should called initSession() at '/ws' ", func() {
		_, err := http.Get("http://localhost:" + port + "/ws")
		Ω(err).ShouldNot(HaveOccurred())

		Ω(fake.InitSessionCallCount()).To(Equal(1))
	})

	//test the handlers in handlers_test
	PIt("should accept websocket at path /ws", func() {
		origin := "http://127.0.0.1/"
		wsUrl := "ws://localhost:" + port + "/ws"

		//disabled gorilla websocket's origin check (origin == r.Host)
		//in order to accept the ws connection
		_, err := websocket.Dial(wsUrl, "ws", origin)
		Ω(err).ShouldNot(HaveOccurred())
	})

	PIt("should path '/' should redirect to '/ui", func() {
		res, err := http.Get("http://localhost:" + port)
		Ω(err).ShouldNot(HaveOccurred())

		Ω(res.Request.URL.String()).Should(Equal("http://localhost:" + port + "/ui/"))
	})

	PIt("should place uploaded files in user dir", func() {

	})

	PIt("should warn about invalid cf commands'", func() {

	})

})
