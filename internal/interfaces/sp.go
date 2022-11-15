package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetTagRepository() ITagRepository
}
