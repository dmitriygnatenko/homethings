package sp

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	authService "git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
	dbService "git.dmitriygnatenko.ru/dima/homethings/internal/services/db"
	envService "git.dmitriygnatenko.ru/dima/homethings/internal/services/env"
	mailerService "git.dmitriygnatenko.ru/dima/homethings/internal/services/mailer"
)

type ServiceProvider struct {
	env                  interfaces.Env
	auth                 interfaces.Auth
	mailer               interfaces.Mailer
	placeRepository      interfaces.PlaceRepository
	thingRepository      interfaces.ThingRepository
	placeThingRepository interfaces.PlaceThingRepository
	placeImageRepository interfaces.PlaceImageRepository
	thingImageRepository interfaces.ThingImageRepository
	userRepository       interfaces.UserRepository
	fileRepository       interfaces.FileRepository
}

func Init() (interfaces.ServiceProvider, error) {
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
	sp.userRepository = repositories.InitUserRepository(db)
	sp.fileRepository = repositories.InitFileRepository()

	return sp, nil
}

func (sp *ServiceProvider) GetEnvService() interfaces.Env {
	return sp.env
}

func (sp *ServiceProvider) GetAuthService() interfaces.Auth {
	return sp.auth
}

func (sp *ServiceProvider) GetMailerService() interfaces.Mailer {
	return sp.mailer
}

func (sp *ServiceProvider) GetPlaceRepository() interfaces.PlaceRepository {
	return sp.placeRepository
}

func (sp *ServiceProvider) GetThingRepository() interfaces.ThingRepository {
	return sp.thingRepository
}

func (sp *ServiceProvider) GetPlaceThingRepository() interfaces.PlaceThingRepository {
	return sp.placeThingRepository
}

func (sp *ServiceProvider) GetPlaceImageRepository() interfaces.PlaceImageRepository {
	return sp.placeImageRepository
}

func (sp *ServiceProvider) GetThingImageRepository() interfaces.ThingImageRepository {
	return sp.thingImageRepository
}

func (sp *ServiceProvider) GetUserRepository() interfaces.UserRepository {
	return sp.userRepository
}

func (sp *ServiceProvider) GetFileRepository() interfaces.FileRepository {
	return sp.fileRepository
}

func InitMock(deps ...interface{}) interfaces.ServiceProvider {
	sp := ServiceProvider{}

	for _, d := range deps {
		switch s := d.(type) {
		case interfaces.Env:
			sp.env = s
		case interfaces.Auth:
			sp.auth = s
		case interfaces.PlaceThingRepository:
			sp.placeThingRepository = s
		case interfaces.ThingRepository:
			sp.thingRepository = s
		case interfaces.PlaceRepository:
			sp.placeRepository = s
		case interfaces.PlaceImageRepository:
			sp.placeImageRepository = s
		case interfaces.ThingImageRepository:
			sp.thingImageRepository = s
		case interfaces.UserRepository:
			sp.userRepository = s
		case interfaces.FileRepository:
			sp.fileRepository = s
		}
	}

	return &sp
}
