package commands

import (
	"fmt"

	"github.com/cloudfoundry-community/cfplayground/cf"
	"github.com/cloudfoundry-community/cfplayground/config"
)

func Admin_CreateNewUser(basePath string, token string, configs *config.Config) (*config.Config, error) {
	//testing new cf admin
	cfAdmin := cf.NewCfAdmin(basePath, configs)
	adminErr := cfAdmin.Login()
	if adminErr != nil {
		return nil, fmt.Errorf("Error logging in as admin: %s", adminErr.Error())
	}

	//create user
	if adminErr == nil {
		adminErr = cfAdmin.CreateUser(token)
		if adminErr != nil {
			return nil, fmt.Errorf("Error creating new user: %s", adminErr.Error())
		}
	}

	//create space
	if adminErr == nil {
		adminErr = cfAdmin.CreateSpace(token)
		if adminErr != nil {
			return nil, fmt.Errorf("Error creating new space: %s", adminErr.Error())
		}
	}

	//assign user to space
	if adminErr == nil {
		adminErr = cfAdmin.AssignUserRole(token, token)
		if adminErr != nil {
			return nil, fmt.Errorf("Error assigning new user role: %s", adminErr.Error())
		}
	}

	// send new user info if all actions are successful
	configs.Server.Login = token
	configs.Server.Pass = "cfplayground"
	configs.Server.Space = token

	return configs, nil
}
