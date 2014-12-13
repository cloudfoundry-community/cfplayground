package users

import (
	"fmt"
	"path"
	"time"

	"github.com/cloudfoundry-community/cfplayground/cf"
	"github.com/cloudfoundry-community/cfplayground/tutorials"
	"github.com/cloudfoundry-community/cfplayground/websocket"

	"github.com/nu7hatch/gouuid"
)

type UniqueUser struct {
	Token       string
	CF          cf.CLI
	Pipe        *websocket.Pipe
	LastConnect time.Time
	Tutorials   *tutorials.TutorialsInfo
	// username    string
	// password    string
}

var userList map[string]UniqueUser

func init() {
	userList = make(map[string]UniqueUser)
}

func New(basePath, token string, cli cf.CLI, pipe *websocket.Pipe) UniqueUser {
	var newUser = UniqueUser{
		token,
		cli,
		pipe,
		time.Now(),
		tutorials.New(path.Join(basePath, "tutorials/courses/")),
		// "",
		// "",
	}
	userList[newUser.Token] = newUser
	return newUser
}

func RestoreUser(token string, cli cf.CLI, pipe *websocket.Pipe) (UniqueUser, error) {
	if user, ok := userList[token]; ok {
		user.CF = cli
		user.Pipe = pipe
		user.LastConnect = time.Now()
		return user, nil
	}

	return UniqueUser{}, fmt.Errorf("User %s not found", token)
}

func List(token string) UniqueUser {
	return userList[token]
}

func GenerateToken() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return guid.String()
}
