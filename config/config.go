package config

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type GlobalExec struct {
	Connection *ssh.Client
	Config     Config
}

type Config struct {
	Repo       string `yaml:"repo"`
	Entrypoint string `yaml:"entrypoint"`
	HomeDir    string `yaml:"homedir"`
	HostAlias  string `yaml:"hostalias"`
	Mode       string `yaml:"mode"`
}

func ReadConfigFile(configFilePath string) (*Config, error) {
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
