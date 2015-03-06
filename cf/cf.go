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
	Marketplace() error
	Services() error
	Files(string, string) error
	BindService(string, string) error
	UnBindService(string, string) error
	App(string) error
	Env(string) error
	Push(string, string, string, int, string) error
	Delete(string) error
	Output(websocket.Message)
	EnvVar() string
	Status() StatusType
	Help(string) error
	Scale(string, string, int, string) error
	Logs(string, bool) error
	Start(string) error
	Stop(string) error
	Restage(string) error
	Restart(string) error
	DeleteService(string) error
	Domains() error
	CreateUserProvidedService(string, string) error
	MapRoute(string, string, string) error
	UnMapRoute(string, string, string) error
}

type CF struct {
	envVar  string
	cfPath  string
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

	err = CopyDir(path.Join(basePath, "assets", "defaultapp"), path.Join(userFolder, "defaultapp"))
	if err != nil {
		panic("Cannot copy default CF App to user directory: " + path.Join(userFolder, "app"))
	}

	absPath, _ := filepath.Abs(userFolder)
	absCfPath, _ := filepath.Abs(path.Join(basePath, "assets/cf/"))
	configs, err := config.New("./config/config.json")

	if err != nil {
		panic("Failed to read config file " + err.Error())
	}

	return &CF{absPath, absCfPath, StatusType{WAITING, ""}, out, in, prompt, configs}
}

func (c *CF) ShowCommand(args ... string) error {
	cmd := exec.Command(path.Join(c.cfPath, "pcf"), args...)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{args[0], "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf apps: ", err)
	}

	cmd.Wait()
	return nil
}
func (c *CF) Login() error {
	c.setStatus(COMMAND, "login")
	defer c.resetStatus()
	fmt.Println("login1: ", c.status.Job, " ", c.status.Detail)
	cmd := exec.Command(path.Join(c.cfPath, "pcf"), "login", "-a", c.configs.Server.Url, "-u", c.configs.Server.Login, "-p", c.configs.Server.Pass, "-o", c.configs.Server.Org, "-s", c.configs.Server.Space, "--skip-ssl-validation")
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"login", "stdout", c.out}

	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf login: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Apps() error {
	return c.ShowCommand("apps")
}

func (c *CF) Services() error {
	return c.ShowCommand("services")
}
func (c *CF) Files(appName string, pathFile string) error {
	if pathFile != "" {
		return c.ShowCommand("files", appName, pathFile)
	}
	return c.ShowCommand("files", appName)
}
func (c *CF) BindService(appName string, serviceName string) error {
	return c.ShowCommand("bind-service", appName, serviceName)
}
func (c *CF) UnBindService(appName string, serviceName string) error {
	return c.ShowCommand("unbind-service", appName, serviceName)
}
func (c *CF) Buildpacks() error {
	return c.ShowCommand("buildpacks")
}
func (c *CF) Domains() error {
	return c.ShowCommand("domains")
}
func (c *CF) Marketplace() error {
	return c.ShowCommand("marketplace")
}
func (c *CF) CreateUserProvidedService(serviceName string, credentials string) error {
	return c.ShowCommand("cups", serviceName, "-p", credentials)
}
func (c *CF) Help(commandName string) error {
	if commandName != "" {
		return c.ShowCommand("help", commandName)
	}
	return c.ShowCommand("help")
}

func (c *CF) App(appName string) error {
	return c.ShowCommand("app", appName)
}
func (c *CF) Start(appName string) error {
	return c.ShowCommand("start", appName)
}
func (c *CF) Stop(appName string) error {
	return c.ShowCommand("stop", appName)
}
func (c *CF) Restage(appName string) error {
	return c.ShowCommand("restage", appName)
}
func (c *CF) Restart(appName string) error {
	return c.ShowCommand("restart", appName)
}
func (c *CF) Env(appName string) error {
	return c.ShowCommand("env", appName)
}
func (c *CF) DeleteService(serviceName string) error {
	return c.ShowCommand("delete-service", serviceName, "-f")
}
func (c *CF) MapRoute(appName string, domain string, hostname string) error {
	return c.ShowCommand("map-route", appName, domain, "-n", hostname)
}
func (c *CF) UnMapRoute(appName string, domain string, hostname string) error {
	return c.ShowCommand("unmap-route", appName, domain, "-n", hostname)
}
func (c *CF) Logs(appName string, recent bool) error {
	c.setStatus(INPUT, "scale")
	defer c.resetStatus()
	var cmd *exec.Cmd;
	if (recent) {
		cmd = exec.Command(path.Join(c.cfPath, "pcf"), "logs", appName, "--recent")
	}else {
		cmd = exec.Command(path.Join(c.cfPath, "pcf"), "logs", appName)
		//when tail we need to send a signal to the cli to exit the tailing.
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
	cmd := exec.Command(path.Join(c.cfPath, "pcf"), "push", appName, "-p", pathApp, "-m", memory, "-i", strconv.Itoa(numberInstance), "-k", diskLimit)
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
	cmd := exec.Command(path.Join(c.cfPath, "pcf"), "scale", appName, "-m", memory, "-i", strconv.Itoa(numberInstance), "-k", diskLimit)
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
	cmd := exec.Command(path.Join(c.cfPath, "pcf"), "delete", appName)
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
