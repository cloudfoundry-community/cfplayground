package server

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/cloudfoundry-community/cfplayground/cf"
	. "github.com/cloudfoundry-community/cfplayground/commands"
	"github.com/cloudfoundry-community/cfplayground/users"
	"github.com/cloudfoundry-community/cfplayground/websocket"

	"github.com/gorilla/mux"
)

var Port string = "8080"

type Server interface {
	Serve(h ServerHandlers)
}

func Serve(h ServerHandlers) {
	RegisterHandler(h)
	bind()
}

func RegisterHandler(h ServerHandlers) {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(h.RedirectBase)
	r.Methods("POST").Path("/upload/{token}").HandlerFunc(h.UploadHandler)

	http.HandleFunc("/ws", h.InitSession)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir(path.Join(h.BasePath(), "ui")))))
	http.Handle("/", r)
}

func bind() {
	fmt.Printf("Starting web ui on http://localhost:%s", Port)
	if err := listenAndServe(":" + Port); err != nil {
		panic(err)
	}
}

var listenAndServe = func(bind string) error {
	return http.ListenAndServe(bind, nil)
}

func getConsoleInput(user *users.UniqueUser) {
	var inputMsg []byte
	for {
		inputMsg = <-user.Pipe.In
		fmt.Println("pipe msg: ", string(inputMsg), " / jobs:", user.CF.Status().Job)
		if user.CF.Status().Job == cf.WAITING {
			go processConsoleInput(user, string(inputMsg))
		} else if user.CF.Status().Job == cf.INPUT {
			user.Pipe.Prompt <- inputMsg //pass input to command
		} else {
			user.CF.Output(websocket.Message{"echo", "warning", "Please wait till current command finish processing."})
		}
	}
}

func processConsoleInput(user *users.UniqueUser, message string) {
	if message == "cf apps" {
		CfApps(user)
	} else if strings.HasPrefix(message, "cf app ") {
		CfApp(user, strings.Trim(message[7:], " "))
	} else if message == "cf push" {
		CfPush(user)
	} else if strings.HasPrefix(message, "cf delete ") {
		CfDelete(user, strings.Trim(message[10:], " "))
	} else if strings.HasPrefix(message, "[course]") {
		RunCourse(user, message[8:])
	} else {
		user.CF.Output(websocket.Message{"echo", "warning", message + " is not a valid command"})
	}

	if user.Tutorials.InProgress() && !strings.HasPrefix(message, "[course]") {
		ProgressCourse(user, message)
	}
}
