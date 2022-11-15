package sp

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	dbService "git.dmitriygnatenko.ru/dima/homethings/internal/services/db"
	envService "git.dmitriygnatenko.ru/dima/homethings/internal/services/env"
)

type ServiceProvider struct {
	env           interfaces.IEnv
	tagRepository interfaces.ITagRepository
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
	sp.tagRepository = repositories.InitTagRepository(db)

	return sp, nil
}

func (sp *ServiceProvider) GetEnvService() interfaces.IEnv {
	return sp.env
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
		}
	}

	return &sp
}
