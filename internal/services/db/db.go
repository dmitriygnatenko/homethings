package db

import (
	"database/sql"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
)

func Init(env interfaces.IEnv) (*sql.DB, error) {
	dataSource := env.GetDBUser() + ":" + env.GetDBPassword() +
		"@tcp(" + env.GetDBHost() + ":" + env.GetDBPort() + ")/" + env.GetDBName()

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
