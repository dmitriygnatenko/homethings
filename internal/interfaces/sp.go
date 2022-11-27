package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetPlaceRepository() IPlaceRepository
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetTagRepository() ITagRepository
}
