package flag

import (
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

	return res, nil
}

func (f *flags) IsCommand() bool {
	return f.help || f.action != "" || f.config == ""
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
