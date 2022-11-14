package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetArticleRepository() IArticleRepository
	GetTagRepository() ITagRepository
	GetArticleTagRepository() IArticleTagRepository
}
