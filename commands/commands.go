package commands

import (
	"github.com/cloudfoundry-community/cfplayground/users"
	"github.com/cloudfoundry-community/cfplayground/websocket"
)

func CfLogin(user *users.UniqueUser) {
	err := user.CF.Login()
	if err != nil {

	}
}

func CfApps(user *users.UniqueUser) {
	user.CF.Output(websocket.Message{"echo", "input", "cf apps"})
	err := user.CF.Apps()
	if err != nil {
	}
}

func CfApp(user *users.UniqueUser, appName string) {
	user.Pipe.Out <- &websocket.Message{"echo", "input", "cf app " + appName}
	err := user.CF.App(appName)
	if err != nil {

	} else {

	}
}

func CfDelete(user *users.UniqueUser, appName string) {
	user.Pipe.Out <- &websocket.Message{"echo", "input", "cf delete " + appName}
	err := user.CF.Delete(appName)
	if err != nil {
	}
}

func CfPush(user *users.UniqueUser, appName string) {
	user.Pipe.Out <- &websocket.Message{"echo", "input", "cf push " + appName}
	err := user.CF.Push(appName)
	if err != nil {
	}
}

func RunCourse(user *users.UniqueUser, cName string) {
	if user.Tutorials.InProgress() {
		user.Pipe.Out <- &websocket.Message{"echo", "warning", "Another Course is currently in progress, you can choose to terminal this course in the dropdown menu"}
		return
	}
	instruction, step := user.Tutorials.StartCourse(cName)
	if instruction == "" {
		user.Pipe.Out <- &websocket.Message{"echo", "warning", "Course is not yet available"}
	} else {
		user.Pipe.Out <- &websocket.Message{"course", cName + " - Step " + step, instruction}
		user.Pipe.Out <- &websocket.Message{"echo", "system", "Please follow the tutorial instruction ..."}
	}
}

func ProgressCourse(user *users.UniqueUser, commandDone string) {
	instruction, cName, step, err := user.Tutorials.ProgressCourse(commandDone)
	if err != nil {
		user.Pipe.Out <- &websocket.Message{"echo", "warning", err.Error()}
	} else if instruction == "" {
		user.Pipe.Out <- &websocket.Message{"echo", "warning", "There is a problem with the next step of this course, please report to admin"}
	} else {
		user.Pipe.Out <- &websocket.Message{"course", cName + " - Step " + step, instruction}
	}
}
