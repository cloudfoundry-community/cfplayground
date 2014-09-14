package websocket_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "github.com/cloudfoundry-community/cfplayground/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	wsClient "code.google.com/p/go.net/websocket"
)

//temporary disabling this test due to gorilla's "origin not allowed" limiation
var _ = XDescribe("Websocket", func() {
	var (
		ts   *httptest.Server
		pipe *Pipe
		ws   net.Conn
		err  error
	)

	BeforeEach(func() {

		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pipe, err = New(w, r)
			Ω(err).ShouldNot(HaveOccurred())
		}))

		parseUrl, err := url.Parse(ts.URL)
		Ω(err).ShouldNot(HaveOccurred())

		origin := "http://127.0.0.1/"
		wsUrl := "ws://" + parseUrl.Host

		//disabled gorilla websocket's origin check (origin == r.Host)
		//origin check blocks the ws connection
		ws, err = wsClient.Dial(wsUrl, "ws", origin)

		Ω(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		ts.Close()
	})

	It("should receive message in the 'In' channel from client", func() {
		ws.Write([]byte("hello there"))

		msg := <-pipe.In
		Ω(string(msg)).To(Equal("hello there"))
	})

	It("sends message to client when 'Out' channel recieve input", func() {
		msg := &Message{"testCmd", "test", "This is a test"}
		pipe.Out <- msg

		var b = make([]byte, 128)
		n, readErr := ws.Read(b)
		Ω(readErr).ShouldNot(HaveOccurred())

		fmt.Println("n: ", n)
		fmt.Println("msg: ", string(b))
	})

})
