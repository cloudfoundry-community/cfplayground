package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Server ServerInfo `json:"server"`
}

type ServerInfo struct {
	Url   string `json:"url"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
	Org   string `json:"org"`
	Space string `json:"space"`
	SkipSSLValidation bool `json:"skip-ssl-validation"` // defaults to false if not present
}

func New(filePath string) (*Config, error) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no such file found: %s", filePath)
	}
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	var configs Config
	json.Unmarshal(b, &configs)

	if configs.Server.Url == "" {
		return nil, fmt.Errorf("Error: server url is missing")
	} else if configs.Server.Login == "" {
		return nil, fmt.Errorf("Error: server login is missing")
	} else if configs.Server.Pass == "" {
		return nil, fmt.Errorf("Error: server password is missing")
	} else if configs.Server.Org == "" {
		return nil, fmt.Errorf("Error: server organization is missing")
	} else if configs.Server.Space == "" {
		return nil, fmt.Errorf("Error: server space is missing")
	}
    
	return &configs, nil
}
