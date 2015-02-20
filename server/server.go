/**
- courses obj
- left is console and option menu
- right is course narrating and course menu
*/

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
var Url string = "localhost"
var WebSocketPort string = "8080"
var mapsCommands = map[string]interface{}{
	"apps": CfApps,
	"a": CfApps,
	"push": CfPush,
	"scale": CfScale,
	"buildpacks": CfBuildpacks,
	"delete": CfDelete,
	"d": CfDelete,
	"logs": CfLogs,
	"help": CfHelp}

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
	http.HandleFunc("/ui/js/env.js", envHandler)
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

func envHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `var wsport = %s;var wsIP = "%s";`, WebSocketPort, Url)
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

	if strings.Fields(message)[0] == "cf" {
		funcName := strings.Fields(message)[1];
		if _, ok := mapsCommands[funcName]; ok {
			err := mapsCommands[funcName].(func(*users.UniqueUser, string) error)(user, message)
			if err == nil {
				return
			}
		}
	}else if strings.HasPrefix(message, "[course]") {
		RunCourse(user, message[8:])
		return
	}

	if user.Tutorials.InProgress() && !strings.HasPrefix(message, "[course]") {
		ProgressCourse(user, message)
		return
	}
	splittedMessage := strings.Fields(message)
	if splittedMessage[0] == "cf" && len(splittedMessage) >= 2 {
		user.CF.Output(websocket.Message{"echo", "warning", message + " is not a valid command, running cf help for 'cf " + splittedMessage[1] + "'"})
		CfHelp(user, "cf help "+splittedMessage[1])
	}else if splittedMessage[0] == "cf" {
		user.CF.Output(websocket.Message{"echo", "warning", message + " is not a valid command, running cf help"})
		CfHelp(user, "cf help")
	}else {
		user.CF.Output(websocket.Message{"echo", "warning", message + " is not a valid command"})
	}

}
