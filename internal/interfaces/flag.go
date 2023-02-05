package interfaces

type IFlag interface {
	IsRunCommand() bool
	GetConfig() string
	GetAction() string
	GetUsername() string
	GetPassword() string
}
