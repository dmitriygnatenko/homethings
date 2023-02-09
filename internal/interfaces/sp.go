package interfaces

type ServiceProvider interface {
	GetEnvService() Env
	GetAuthService() Auth
	GetMailerService() Mailer
	GetPlaceRepository() PlaceRepository
	GetThingRepository() ThingRepository
	GetPlaceThingRepository() PlaceThingRepository
	GetPlaceImageRepository() PlaceImageRepository
	GetThingImageRepository() ThingImageRepository
	GetUserRepository() UserRepository
	GetFileRepository() FileRepository
}
