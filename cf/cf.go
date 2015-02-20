///////////////////////////////////
/*
issue:
	- cf Environemnt Var is not setting correctly in Windows
	- need cf.exe for windows
*/
///////////////////////////////////
package cf

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"math/rand"
	"time"
	"strconv"
	"github.com/cloudfoundry-community/cfplayground/config"
	. "github.com/cloudfoundry-community/cfplayground/copy"
	"github.com/cloudfoundry-community/cfplayground/websocket"
	"syscall"
)

const (
	WAITING = iota //waiting for command
	INPUT          //performing a command and waiting for input from user
	COMMAND        //performing a command, no input from user is expected
)

type CLI interface {
	Login() error
	Apps() error
	Buildpacks() error
	App(string) error
	Push(string, string, string, int, string) error
	Delete(string) error
	Output(websocket.Message)
	EnvVar() string
	Status() StatusType
	Help(string) error
	Scale(string, string, int, string) error
	Logs(string, bool) error
}

type CF struct {
	envVar  string
	status  StatusType
	out     chan *websocket.Message
	in      chan []byte
	prompt  chan []byte
	configs *config.Config
}

type StatusType struct {
	Job    int
	Detail string
}

type msgWriter struct {
	cmd     string
	msgType string
	out     chan *websocket.Message
}

type msgReader struct {
	in chan []byte
}

func (w *msgWriter) Write(b []byte) (n int, err error) {
	w.out <- &websocket.Message{w.cmd, w.msgType, string(b)}
	return len(b), nil
}

func New(token string, out chan *websocket.Message, in chan []byte, prompt chan []byte, basePath string) CLI {
	containerPath := path.Join(basePath, "containers")
	userFolder := path.Join(containerPath, token)

	err := os.MkdirAll(path.Join(userFolder), os.ModePerm)
	if err != nil {
		panic("Cannot create user directory: " + path.Join(containerPath, token))
	}

	err = CopyFile(path.Join(basePath, "assets/cf/", "pcf"), path.Join(userFolder, "pcf"))
	if err != nil {
		panic("Cannot copy cf binary to user directory: " + path.Join(userFolder))
	}

	err = CopyDir(path.Join(basePath, "assets", "defaultapp"), path.Join(userFolder, "defaultapp"))
	if err != nil {
		panic("Cannot copy default CF App to user directory: " + path.Join(userFolder, "app"))
	}

	absPath, _ := filepath.Abs(userFolder)
	configs, err := config.New("./config/config.json")

	if err != nil {
		panic("Failed to read config file " + err.Error())
	}

	return &CF{absPath, StatusType{WAITING, ""}, out, in, prompt, configs}
}

func (c *CF) Login() error {
	c.setStatus(COMMAND, "login")
	defer c.resetStatus()
	fmt.Println("login1: ", c.status.Job, " ", c.status.Detail)
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "login", "-a", c.configs.Server.Url, "-u", c.configs.Server.Login, "-p", c.configs.Server.Pass, "-o", c.configs.Server.Org, "-s", c.configs.Server.Space, "--skip-ssl-validation")
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"login", "stdout", c.out}

	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf login: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Apps() error {
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "apps")
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"apps", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf apps: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Buildpacks() error {
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "buildpacks")
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"buildpacks", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf buildpacks: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Help(commandName string) error {
	var cmd *exec.Cmd;
	if commandName != "" {
		cmd = exec.Command(path.Join(c.envVar, "pcf"), "help", commandName)
	}else {
		cmd = exec.Command(path.Join(c.envVar, "pcf"), "help")
	}

	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"help", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf help: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) App(appName string) error {
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "app", appName)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"app", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf apps: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Logs(appName string, recent bool) error {
	c.setStatus(INPUT, "scale")
	defer c.resetStatus()
	var cmd *exec.Cmd;
	if (recent) {
		cmd = exec.Command(path.Join(c.envVar, "pcf"), "logs", appName, "--recent")
	}else {
		cmd = exec.Command(path.Join(c.envVar, "pcf"), "logs", appName)

		var msg []byte
		go func() {
			msgW := &msgWriter{"scale", "stdout", c.out}
			for {
				msg = <-c.prompt
				cmd.Process.Signal(syscall.SIGINT)
				msgW.Write(msg)
			}
		}()
	}
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"logs", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error running cf logs: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Push(appName string, pathApp string, memory string, numberInstance int, diskLimit string) error {
	rand.Seed(time.Now().UnixNano())
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "push", appName, "-p", pathApp, "-m", memory, "-i", strconv.Itoa(numberInstance), "-k", diskLimit)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Dir = c.envVar
	cmd.Stdout = &msgWriter{"push", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf push: ", err)
	}
	cmd.Wait()
	return nil
}

func (c *CF) Scale(appName string, memory string, numberInstance int, diskLimit string) error {
	c.setStatus(INPUT, "scale")
	defer c.resetStatus()
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "scale", appName, "-m", memory, "-i", strconv.Itoa(numberInstance), "-k", diskLimit)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")

	stdin, _ := cmd.StdinPipe()

	cmd.Stdout = &msgWriter{"scale", "stdout", c.out}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error running cf scale: ", err)
	}

	var msg []byte
	go func() {
		msgW := &msgWriter{"scale", "stdout", c.out}
		for {
			msg = <-c.prompt
			io.WriteString(stdin, string(msg)+"\n")
			msgW.Write(msg)
		}
	}()

	cmd.Wait()
	return nil
}

func (c *CF) Delete(appName string) error {
	c.setStatus(INPUT, "delete")
	defer c.resetStatus()
	cmd := exec.Command(path.Join(c.envVar, "pcf"), "delete", appName)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")

	stdin, _ := cmd.StdinPipe()
	cmd.Stdout = &msgWriter{"delete", "stdout", c.out}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error in cf delete: ", err)
	}

	var msg []byte
	go func() {
		msgW := &msgWriter{"delete", "stdout", c.out}
		for {
			msg = <-c.prompt
			io.WriteString(stdin, string(msg)+"\n")
			msgW.Write(msg)
		}
	}()

	cmd.Wait()
	return nil
}

func (c *CF) Output(msg websocket.Message) {
	c.out <- &msg
}

func (c *CF) EnvVar() string {
	return c.envVar
}

func (c *CF) Status() StatusType {
	return c.status
}

func (c *CF) setStatus(job int, detail string) {
	c.status.Job = job
	c.status.Detail = detail
}

func (c *CF) resetStatus() {
	c.status.Job = WAITING
	c.status.Detail = ""
}
