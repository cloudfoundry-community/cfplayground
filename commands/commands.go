package commands

import (
	"github.com/cloudfoundry-community/cfplayground/users"
	"github.com/cloudfoundry-community/cfplayground/websocket"
	"strings"
	"flag"
	"fmt"
)

func CfLogin(user *users.UniqueUser) error {
	//scan, err := user.CF.Login()
	//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[start]"))
	err := user.CF.Login()
	if err != nil {
		return err
	} else {
		// for scan.Scan() {
		// 	user.Pipe.Conn.WriteMessage(websocket.TextMessage, scan.Bytes())
		// }
		//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[done!]"))
	}
	return nil
}

func CfApps(user *users.UniqueUser, message string) error {
	user.CF.Output(websocket.Message{"echo", "input", message})
	err := user.CF.Apps()
	if err != nil {
		return err
	}
	return nil
}

func CfApp(user *users.UniqueUser, message string) error {
	user.Pipe.Out <- &websocket.Message{"echo", "input", message}

	if len(strings.Fields(message)) < 3 {
		return fmt.Errorf("command not valid missing appname")
	}

	err := user.CF.App(strings.Fields(message)[2])
	if err != nil {
		return err;
	}
	return nil
}

func CfHelp(user *users.UniqueUser, message string) error {

	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	var err error;
	if (len(strings.Fields(message)) == 3) {
		err = user.CF.Help(strings.Fields(message)[2])
	}else {
		err = user.CF.Help("")
	}

	if err != nil {
		return err
	}
	return nil
}
func CfLogs(user *users.UniqueUser, message string) error {

	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	if len(strings.Fields(message)) < 3 {
		return fmt.Errorf("command not valid missing appname")
	}
	var err error;

	if (len(strings.Fields(message)) == 4 && strings.Fields(message)[3] == "--recent") {
		err = user.CF.Logs(strings.Fields(message)[2], true)
	}else {
		err = user.CF.Logs(strings.Fields(message)[2], false)
	}

	if err != nil {
		return err;
	}
	return nil
}
func CfDelete(user *users.UniqueUser, message string) error {
	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	//scan, err := user.CF.Login()
	//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[start]"))
	if len(strings.Fields(message)) < 3 {
		return fmt.Errorf("command not valid missing appname")
	}
	err := user.CF.Delete(strings.Fields(message)[2])
	if err != nil {
		return err
	} else {
		// for scan.Scan() {
		// 	user.Pipe.Conn.WriteMessage(websocket.TextMessage, scan.Bytes())
		// }
		//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[done!]"))
	}
	return nil
}

func CfPush(user *users.UniqueUser, message string) error {
	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	var CommandLine = flag.NewFlagSet("push", flag.ContinueOnError)
	var pathPush string
	var memory string
	var numberInstance int
	var diskLimit string
	if len(strings.Fields(message)) < 3 {
		return fmt.Errorf("command not valid missing appname")
	}
	CommandLine.StringVar(&pathPush, "p", "", "path")
	CommandLine.StringVar(&memory, "m", "1G", "memory")
	CommandLine.IntVar(&numberInstance, "i", 1, "number of instances")
	CommandLine.StringVar(&diskLimit, "k", "1G", "disk limit")
	CommandLine.Parse(strings.Fields(message)[3:])

	if pathPush == "" {
		return fmt.Errorf("command not valid missing path (option -p)")
	}

	//scan, err := user.CF.Login()
	//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[start]"))
	err := user.CF.Push(strings.Fields(message)[2], pathPush, memory, numberInstance, diskLimit)
	if err != nil {
		return err
	}
	return nil
}

func CfScale(user *users.UniqueUser, message string) error {
	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	var CommandLine = flag.NewFlagSet("scale", flag.ContinueOnError)
	var memory string
	var numberInstance int
	var diskLimit string
	if len(strings.Fields(message)) < 3 {
		return fmt.Errorf("command not valid missing appname")
	}
	CommandLine.StringVar(&memory, "m", "1G", "memory")
	CommandLine.IntVar(&numberInstance, "i", 1, "number of instances")
	CommandLine.StringVar(&diskLimit, "k", "1G", "disk limit")
	CommandLine.Parse(strings.Fields(message)[3:])
	//scan, err := user.CF.Login()
	//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[start]"))
	err := user.CF.Scale(strings.Fields(message)[2], memory, numberInstance, diskLimit)
	if err != nil {
		return err
	}
	return nil
}

func CfBuildpacks(user *users.UniqueUser, message string) error {
	user.Pipe.Out <- &websocket.Message{"echo", "input", message}
	//scan, err := user.CF.Login()
	//user.Pipe.Conn.WriteMessage(websocket.TextMessage, []byte("[start]"))
	err := user.CF.Buildpacks()
	if err != nil {
		return err
	}
	return nil
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
