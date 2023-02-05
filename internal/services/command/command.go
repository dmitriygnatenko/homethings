package command

import (
	"errors"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/flag"
)

type command struct {
	flags interfaces.IFlag
}

func Init(flags interfaces.IFlag) (interfaces.ICommand, error) {
	if flags.GetAction() != "" {
		switch flags.GetAction() {
		case flag.ActionTypeAddUser, flag.ActionTypeUpdateUser:
			if flags.GetUsername() == "" {
				return nil, errors.New("parameter *username* must be filled for the add/update user action")
			}
			if flags.GetPassword() == "" {
				return nil, errors.New("parameter *password* must be filled for the add/update user action")
			}
		case flag.ActionTypeDeleteUser:
			if flags.GetUsername() == "" {
				return nil, errors.New("parameter *username* must be filled for the delete user action")
			}
		default:
			return nil, errors.New("incorrect *action* parameter")
		}
	}

	return &command{
		flags: flags,
	}, nil
}

func (c *command) Run() {
	switch c.flags.GetAction() {
	case flag.ActionTypeAddUser:
		c.addUser(c.flags.GetUsername(), c.flags.GetPassword())
	case flag.ActionTypeUpdateUser:
		c.updateUser(c.flags.GetUsername(), c.flags.GetPassword())
	case flag.ActionTypeDeleteUser:
		c.deleteUser(c.flags.GetUsername())
	}
}

func (c *command) addUser(username string, password string) {

}

func (c *command) updateUser(username string, password string) {

}

func (c *command) deleteUser(username string) {

}
