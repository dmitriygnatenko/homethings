package flag

import (
	"errors"
	"flag"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
)

const (
	CommandFlagConfig   = "config"
	CommandFlagHelp     = "help"
	CommandFlagAction   = "action"
	CommandFlagUsername = "username"
	CommandFlagPassword = "password"

	ActionTypeAddUser    = "add-user"
	ActionTypeUpdateUser = "update-user"
	ActionTypeDeleteUser = "delete-user"
)

type flags struct {
	config   string
	help     bool
	action   string
	username string
	password string
}

func Init() (interfaces.IFlag, error) {
	res := &flags{}

	flag.StringVar(&res.config, CommandFlagConfig, "", "Path to config file")
	flag.BoolVar(&res.help, CommandFlagHelp, false, "Display help info")
	flag.StringVar(&res.action, CommandFlagAction, "", "Action type")
	flag.StringVar(&res.username, CommandFlagUsername, "", "Username")
	flag.StringVar(&res.password, CommandFlagPassword, "", "Password")
	flag.Parse()

	return res, res.validate()
}

func (f *flags) validate() error {
	var err error

	if f.config == "" {
		return errors.New("parameter *config* must be filled")
	}

	switch f.action {
	case ActionTypeAddUser, ActionTypeUpdateUser:
		if f.username == "" {
			return errors.New("parameter *username* must be filled for the add/update user action")
		}
		if f.password == "" {
			return errors.New("parameter *password* must be filled for the add/update user action")
		}
	case ActionTypeDeleteUser:
		if f.username == "" {
			return errors.New("parameter *username* must be filled for the delete user action")
		}
	default:
		f.action = ""
	}

	return err
}

func (f *flags) GetConfig() string {
	return f.config
}

func (f *flags) IsHelp() bool {
	return f.help
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
