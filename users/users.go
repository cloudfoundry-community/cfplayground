package users

import (
	"net/http"
	"path"
	"time"

	"github.com/simonleung8/cfplayground/cf"
	"github.com/simonleung8/cfplayground/tutorials"
	"github.com/simonleung8/cfplayground/websocket"

	"github.com/nu7hatch/gouuid"
)

type UniqueUser struct {
	Token       string
	CF          cf.CLI
	Pipe        *websocket.Pipe
	LastConnect time.Time
	Tutorials   *tutorials.TutorialsInfo
	username    string
	password    string
}

var userList map[string]UniqueUser

func init() {
	userList = make(map[string]UniqueUser)
}

func New(w http.ResponseWriter, r *http.Request, basePath, token string, cli cf.CLI, pipe *websocket.Pipe) UniqueUser {
	var newUser = UniqueUser{
		token,
		cli,
		pipe,
		time.Now(),
		tutorials.New(path.Join(basePath, "tutorials/courses/")),
		"",
		"",
	}
	userList[newUser.Token] = newUser
	return newUser
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
