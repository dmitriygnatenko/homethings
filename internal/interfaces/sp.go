package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetPlaceRepository() IPlaceRepository
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetPlaceImageRepository() IPlaceImageRepository
	GetThingImageRepository() IThingImageRepository
	GetFileRepository() IFileRepository
}
