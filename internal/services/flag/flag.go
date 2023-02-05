package flag

import (
	"errors"
	"flag"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
)

const (
	CommandFlagConfig   = "config"
	CommandFlagAction   = "action"
	CommandFlagUsername = "username"
	CommandFlagPassword = "password"

	ActionTypeAddUser    = "add-user"
	ActionTypeUpdateUser = "update-user"
	ActionTypeDeleteUser = "delete-user"
)

type flags struct {
	config   string
	action   string
	username string
	password string
}

func Init() (interfaces.IFlag, error) {
	res := &flags{}

	flag.StringVar(&res.config, CommandFlagConfig, "", "Path to config file")
	flag.StringVar(&res.action, CommandFlagAction, "", "Action type")
	flag.StringVar(&res.username, CommandFlagUsername, "", "Username")
	flag.StringVar(&res.password, CommandFlagPassword, "", "Password")
	flag.Parse()

	if res.config == "" {
		return nil, errors.New("parameter *config* must be filled")
	}

	return res, nil
}

func (f *flags) IsRunCommand() bool {
	return f.action != ""
}

func (f *flags) GetConfig() string {
	return f.config
}

func (f *flags) GetAction() string {
	return f.action
}

func (f *flags) GetUsername() string {
	return f.username
}

func (f *flags) GetPassword() string {
	return f.password
}
