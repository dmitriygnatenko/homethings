package interfaces

type IFlag interface {
	IsCommand() bool
	GetConfig() string
	IsHelp() bool
	GetAction() string
	GetUsername() string
	GetPassword() string
}
