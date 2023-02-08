package interfaces

type IServiceProvider interface {
	GetEnvService() IEnv
	GetAuthService() IAuth
	GetMailerService() IMailer
	GetPlaceRepository() IPlaceRepository
	GetThingRepository() IThingRepository
	GetPlaceThingRepository() IPlaceThingRepository
	GetPlaceImageRepository() IPlaceImageRepository
	GetThingImageRepository() IThingImageRepository
	GetUserRepository() IUserRepository
	GetFileRepository() IFileRepository
}
