package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/cfplayground/server"
	"github.com/cloudfoundry-community/go-cfenv"
)

func main() {
	var port string
	var url string
	var webSocketPort string
	fmt.Println(os.Getenv("VCAP_APPLICATION"))
	appEnv, enverr := cfenv.Current()
	if enverr != nil {
		port = "8080"
		url = "localhost"
		webSocketPort = "8080"
	} else {
		port = strconv.Itoa(appEnv.Port)
		url = appEnv.ApplicationUri[0]
		webSocketPort = "80"
	}
	server.Port = port
	server.Url = url
	server.WebSocketPort = webSocketPort
	server.Serve(server.NewHandler("./"))
}
