package cf

import (
	"os/exec"
	"path/filepath"

	"github.com/cloudfoundry-community/cfplayground/config"
)

type Cf_Admin interface {
	Login() error
	CreateUser(string) error
	CreateSpace(string) error
	AssignUserRole(string, string) error
}

type CfAdmin struct {
	configs *config.Config
	status  StatusType
}

var basePath string

func NewCfAdmin(basePath string, configs *config.Config) Cf_Admin {
	basePath = basePath
	return &CfAdmin{
		configs: configs,
		status:  StatusType{WAITING, ""},
	}
}

func (c *CfAdmin) Login() error {
	c.setStatus(ADMIN, "login")
	defer c.resetStatus()
	sslValidation := ""
	if c.configs.Server.SkipSSLValidation {
		sslValidation = "--skip-ssl-validation"
	}
	cmd := exec.Command(filepath.Join(basePath, "assets/cf/", "pcf"), "login", "-a", c.configs.Server.Url, "-u",
		c.configs.Server.Login, "-p", c.configs.Server.Pass, "-o", c.configs.Server.Org,
		"-s", c.configs.Server.Space, sslValidation)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *CfAdmin) CreateUser(newUser string) error {
	c.setStatus(ADMIN, "createUser")
	defer c.resetStatus()

	cmd := exec.Command(filepath.Join(basePath, "assets/cf/", "pcf"), "create-user", newUser, "password")

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *CfAdmin) AssignUserRole(user, space string) error {
	c.setStatus(ADMIN, "assignUserRole")
	defer c.resetStatus()

	cmd := exec.Command(filepath.Join(basePath, "assets/cf/", "pcf"), "set-space-role", user, c.configs.Server.Org, space, "SpaceDeveloper")

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *CfAdmin) CreateSpace(newSpace string) error {
	c.setStatus(ADMIN, "createSpace")
	defer c.resetStatus()

	cmd := exec.Command(filepath.Join(basePath, "assets/cf/", "pcf"), "create-space", newSpace)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *CfAdmin) setStatus(job int, detail string) {
	c.status.Job = job
	c.status.Detail = detail
}

func (c *CfAdmin) resetStatus() {
	c.status.Job = WAITING
	c.status.Detail = ""
}
