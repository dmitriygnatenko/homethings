package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetAuthService() IAuth
	GetPlaceRepository() IPlaceRepository
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetPlaceImageRepository() IPlaceImageRepository
	GetThingImageRepository() IThingImageRepository
	GetUserRepository() IUserRepository
	GetFileRepository() IFileRepository
}
