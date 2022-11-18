package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetTagRepository() ITagRepository
}
