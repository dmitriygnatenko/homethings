package interfaces

type IServiceProvider interface {
	GetFlagService() IFlag
	GetEnvService() IEnv
	GetCommandService() ICommand
	GetPlaceRepository() IPlaceRepository
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetPlaceImageRepository() IPlaceImageRepository
	GetThingImageRepository() IThingImageRepository
	GetFileRepository() IFileRepository
}
