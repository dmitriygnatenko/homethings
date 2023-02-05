package command

import (
	"context"
	"errors"
	"fmt"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/mappers"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/flag"
	"golang.org/x/crypto/bcrypt"
)

type command struct {
	flags          interfaces.IFlag
	userRepository interfaces.IUserRepository
}

func Init(flags interfaces.IFlag, userRepository interfaces.IUserRepository) (interfaces.ICommand, error) {
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
		flags:          flags,
		userRepository: userRepository,
	}, nil
}

func (c *command) Run() {
	ctx := context.Background()

	switch c.flags.GetAction() {
	case flag.ActionTypeAddUser:
		c.addUser(ctx, c.flags.GetUsername(), c.flags.GetPassword())
	case flag.ActionTypeUpdateUser:
		c.updateUser(ctx, c.flags.GetUsername(), c.flags.GetPassword())
	case flag.ActionTypeDeleteUser:
		c.deleteUser(ctx, c.flags.GetUsername())
	}
}

func (c *command) addUser(ctx context.Context, username string, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	_, err = c.userRepository.Add(ctx, mappers.ConvertToAddUserRequestModel(username, string(hashedPassword)))
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("User successfully added")
}

func (c *command) updateUser(ctx context.Context, username string, password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	err = c.userRepository.Update(ctx, mappers.ConvertToUpdateUserRequestModel(username, string(hashedPassword)))
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("User successfully updated")
}

func (c *command) deleteUser(ctx context.Context, username string) {
	if err := c.userRepository.Delete(ctx, username); err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	fmt.Println("User successfully deleted")
}
