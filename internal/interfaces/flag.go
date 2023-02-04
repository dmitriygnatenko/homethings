package interfaces

type IFlag interface {
	GetConfig() string
	IsHelp() bool
	GetAction() string
	GetUsername() string
	GetPassword() string
}
