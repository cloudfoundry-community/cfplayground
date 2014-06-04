package main

import "github.com/simonleung8/cfplayground/server"

func main() {
	server.Serve(server.NewHandler("./"))
}
