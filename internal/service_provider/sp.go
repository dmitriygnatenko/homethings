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
	tagRepository        interfaces.ITagRepository
}

func Init() (interfaces.IServiceProvider, error) {
	sp := &ServiceProvider{}

	// Init services
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
	sp.tagRepository = repositories.InitTagRepository(db)

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

func (sp *ServiceProvider) GetTagRepository() interfaces.ITagRepository {
	return sp.tagRepository
}

func InitMock(deps ...interface{}) interfaces.IServiceProvider {
	sp := ServiceProvider{}

	for _, d := range deps {
		switch s := d.(type) {
		case interfaces.IEnv:
			sp.env = s
		case interfaces.ITagRepository:
			sp.tagRepository = s
		case interfaces.IPlaceThingRepository:
			sp.placeThingRepository = s
		case interfaces.IThingRepository:
			sp.thingRepository = s
		case interfaces.IPlaceRepository:
			sp.placeRepository = s
		}
	}

	return &sp
}
