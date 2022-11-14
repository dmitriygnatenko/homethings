package interfaces

type IEnv interface {
	GetAppPort() string
	GetDBPort() string
	GetDBHost() string
	GetDBName() string
	GetDBUser() string
	GetDBPassword() string
}
