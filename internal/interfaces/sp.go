package interfaces

type ServiceProvider interface {
	GetEnvService() Env
	GetAuthService() Auth
	GetMailerService() Mailer
	GetPlaceRepository() PlaceRepository
	GetThingRepository() ThingRepository
	GetTagRepository() TagRepository
	GetPlaceThingRepository() PlaceThingRepository
	GetPlaceImageRepository() PlaceImageRepository
	GetThingImageRepository() ThingImageRepository
	GetThingTagRepository() ThingTagRepository
	GetUserRepository() UserRepository
	GetFileRepository() FileRepository
}
