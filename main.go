package main

import "github.com/cloudfoundry-community/cfplayground/server"

func main() {
	server.Serve(server.NewHandler("./"))
}
