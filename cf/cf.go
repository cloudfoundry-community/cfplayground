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
	"path/filepath"

	"github.com/cloudfoundry-community/cfplayground/config"
	. "github.com/cloudfoundry-community/cfplayground/copy"
	"github.com/cloudfoundry-community/cfplayground/utils"
	"github.com/cloudfoundry-community/cfplayground/websocket"
)

type CLI interface {
	Login() error
	Apps() error
	App(string) error
	Push(string) error
	Delete(string) error
	Output(websocket.Message)
	EnvVar() string
	Status() StatusType
}

type CF struct {
	envVar  string
	status  StatusType
	out     chan *websocket.Message
	in      chan []byte
	prompt  chan []byte
	configs *config.Config
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

func NewCli(token string, out chan *websocket.Message, in chan []byte, prompt chan []byte, basePath string, configs *config.Config) CLI {
	containerPath := filepath.Join(basePath, "containers")
	userFolder := filepath.Join(containerPath, token)

	err := os.MkdirAll(filepath.Join(userFolder), os.ModePerm)
	if err != nil {
		panic("Cannot create user directory: " + filepath.Join(containerPath, token))
	}

	err = CopyFile(filepath.Join(basePath, "assets/cf/", "pcf"), filepath.Join(userFolder, "pcf"))
	if err != nil {
		panic("Cannot copy cf binary to user directory: " + filepath.Join(userFolder))
	}

	err = CopyDir(filepath.Join(basePath, "assets", "dora"), filepath.Join(userFolder, "dora"))
	if err != nil {
		panic("Cannot copy default CF App to user directory: " + filepath.Join(userFolder, "app"))
	}

	absPath, _ := filepath.Abs(userFolder)

	return &CF{
		absPath,
		StatusType{WAITING, ""},
		out,
		in,
		prompt,
		configs,
	}
}

func (c *CF) Login() error {
	c.setStatus(COMMAND, "login")
	defer c.resetStatus()
	fmt.Println("login1: ", c.status.Job, " ", c.status.Detail)
	sslValidation := ""
	if c.configs.Server.SkipSSLValidation {
		sslValidation = "--skip-ssl-validation"
	}
	cmd := exec.Command(filepath.Join(c.envVar, "pcf"), "login", "-a", c.configs.Server.Url, "-u",
		c.configs.Server.Login, "-p", c.configs.Server.Pass, "-o", c.configs.Server.Org,
		"-s", c.configs.Server.Space, sslValidation)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"login", "stdout", c.out}

	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf login: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Apps() error {
	cmd := exec.Command(filepath.Join(c.envVar, "pcf"), "apps")
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"apps", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf apps: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) App(appName string) error {
	cmd := exec.Command(filepath.Join(c.envVar, "pcf"), "app", appName)
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Stdout = &msgWriter{"app", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf apps: ", err)
	}

	cmd.Wait()
	return nil
}

func (c *CF) Push(appName string) error {
	var cmd *exec.Cmd

	//push user folder if there is one
	if utils.IsDirExists(filepath.Join(c.envVar, "app")) {
		cmd = exec.Command(filepath.Join(c.envVar, "pcf"), "push", appName, "-p", "app/")
	} else {
		cmd = exec.Command(filepath.Join(c.envVar, "pcf"), "push", appName, "-p", "dora/")
	}
	cmd.Env = append(cmd.Env, "CF_HOME="+c.envVar, "CF_COLOR=true")
	cmd.Dir = c.envVar
	cmd.Stdout = &msgWriter{"push", "stdout", c.out}
	if err := cmd.Start(); err != nil {
		fmt.Errorf("Error running cf push: ", err)
	}
	cmd.Wait()
	return nil
}

func (c *CF) Delete(appName string) error {
	c.setStatus(INPUT, "delete")
	defer c.resetStatus()
	cmd := exec.Command(filepath.Join(c.envVar, "pcf"), "delete", appName)
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
