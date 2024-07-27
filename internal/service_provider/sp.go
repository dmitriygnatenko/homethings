package sp

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	authService "git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
	dbService "git.dmitriygnatenko.ru/dima/homethings/internal/services/db"
	envService "git.dmitriygnatenko.ru/dima/homethings/internal/services/env"
	mailerService "git.dmitriygnatenko.ru/dima/homethings/internal/services/mailer"
)

type ServiceProvider struct {
	env                         *envService.Service
	auth                        *authService.Service
	mailer                      *mailerService.Service
	placeRepository             *repositories.PlaceRepository
	thingRepository             *repositories.ThingRepository
	tagRepository               *repositories.TagRepository
	placeThingRepository        *repositories.PlaceThingRepository
	placeImageRepository        *repositories.PlaceImageRepository
	thingImageRepository        *repositories.ThingImageRepository
	thingTagRepository          *repositories.ThingTagRepository
	userRepository              *repositories.UserRepository
	thingNotificationRepository *repositories.ThingNotificationRepository
	fileRepository              *repositories.FileRepository
}

func Init() (*ServiceProvider, error) {
	sp := &ServiceProvider{}

	env, err := envService.Init()
	if err != nil {
		return nil, err
	}
	sp.env = env

	auth, err := authService.Init(env)
	if err != nil {
		return nil, err
	}
	sp.auth = auth

	mailer, err := mailerService.Init(env)
	if err != nil {
		return nil, err
	}
	sp.mailer = mailer

	db, err := dbService.Init(env)
	if err != nil {
		return nil, err
	}

	// Init repositories
	sp.placeRepository = repositories.InitPlaceRepository(db)
	sp.thingRepository = repositories.InitThingRepository(db)
	sp.placeThingRepository = repositories.InitPlaceThingRepository(db)
	sp.placeImageRepository = repositories.InitPlaceImageRepository(db)
	sp.thingImageRepository = repositories.InitThingImageRepository(db)
	sp.tagRepository = repositories.InitTagRepository(db)
	sp.thingTagRepository = repositories.InitThingTagRepository(db)
	sp.userRepository = repositories.InitUserRepository(db)
	sp.thingNotificationRepository = repositories.InitThingNotificationRepository(db)
	sp.fileRepository = repositories.InitFileRepository()

	return sp, nil
}

func (sp *ServiceProvider) GetEnvService() *envService.Service {
	return sp.env
}

func (sp *ServiceProvider) GetAuthService() *authService.Service {
	return sp.auth
}

func (sp *ServiceProvider) GetMailerService() *mailerService.Service {
	return sp.mailer
}

func (sp *ServiceProvider) GetPlaceRepository() *repositories.PlaceRepository {
	return sp.placeRepository
}

func (sp *ServiceProvider) GetThingRepository() *repositories.ThingRepository {
	return sp.thingRepository
}

func (sp *ServiceProvider) GetTagRepository() *repositories.TagRepository {
	return sp.tagRepository
}

func (sp *ServiceProvider) GetPlaceThingRepository() *repositories.PlaceThingRepository {
	return sp.placeThingRepository
}

func (sp *ServiceProvider) GetPlaceImageRepository() *repositories.PlaceImageRepository {
	return sp.placeImageRepository
}

func (sp *ServiceProvider) GetThingImageRepository() *repositories.ThingImageRepository {
	return sp.thingImageRepository
}

func (sp *ServiceProvider) GetThingTagRepository() *repositories.ThingTagRepository {
	return sp.thingTagRepository
}

func (sp *ServiceProvider) GetThingNotificationRepository() *repositories.ThingNotificationRepository {
	return sp.thingNotificationRepository
}

func (sp *ServiceProvider) GetUserRepository() *repositories.UserRepository {
	return sp.userRepository
}

func (sp *ServiceProvider) GetFileRepository() *repositories.FileRepository {
	return sp.fileRepository
}
