package sp

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	dbService "git.dmitriygnatenko.ru/dima/homethings/internal/services/db"
	envService "git.dmitriygnatenko.ru/dima/homethings/internal/services/env"
)

type ServiceProvider struct {
	env                  interfaces.IEnv
	placeRepository      interfaces.IPlaceRepository
	thingRepository      interfaces.IThingRepository
	placeThingRepository interfaces.IPlaceThingRepository
	placeImageRepository interfaces.IPlaceImageRepository
	thingImageRepository interfaces.IThingImageRepository
	userRepository       interfaces.IUserRepository
	fileRepository       interfaces.IFileRepository
}

func Init() (interfaces.IServiceProvider, error) {
	sp := &ServiceProvider{}

	env, err := envService.Init()
	if err != nil {
		return nil, err
	}
	sp.env = env

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

func (sp *ServiceProvider) GetEnvService() interfaces.IEnv {
	return sp.env
}

func (sp *ServiceProvider) GetPlaceRepository() interfaces.IPlaceRepository {
	return sp.placeRepository
}

func (sp *ServiceProvider) GetThingRepository() interfaces.IThingRepository {
	return sp.thingRepository
}

func (sp *ServiceProvider) GetPlaceThingRepository() interfaces.IPlaceThingRepository {
	return sp.placeThingRepository
}

func (sp *ServiceProvider) GetPlaceImageRepository() interfaces.IPlaceImageRepository {
	return sp.placeImageRepository
}

func (sp *ServiceProvider) GetThingImageRepository() interfaces.IThingImageRepository {
	return sp.thingImageRepository
}

func (sp *ServiceProvider) GetUserRepository() interfaces.IUserRepository {
	return sp.userRepository
}

func (sp *ServiceProvider) GetFileRepository() interfaces.IFileRepository {
	return sp.fileRepository
}

func InitMock(deps ...interface{}) interfaces.IServiceProvider {
	sp := ServiceProvider{}

	for _, d := range deps {
		switch s := d.(type) {
		case interfaces.IEnv:
			sp.env = s
		case interfaces.IPlaceThingRepository:
			sp.placeThingRepository = s
		case interfaces.IThingRepository:
			sp.thingRepository = s
		case interfaces.IPlaceRepository:
			sp.placeRepository = s
		case interfaces.IPlaceImageRepository:
			sp.placeImageRepository = s
		case interfaces.IThingImageRepository:
			sp.thingImageRepository = s
		case interfaces.IUserRepository:
			sp.userRepository = s
		case interfaces.IFileRepository:
			sp.fileRepository = s
		}
	}

	return &sp
}
