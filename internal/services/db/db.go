package db

import (
	"database/sql"
	"time"
)

type Env interface {
	GetAppPort() string
	GetDBHost() string
	GetDBPort() string
	GetDBName() string
	GetDBUser() string
	GetDBPassword() string
	GetDBMaxOpenConns() int
	GetDBMaxIdleConns() int
	GetDBMaxConnLifetime() int
	GetDBMaxIdleConnLifetime() int
}

func Init(env Env) (*sql.DB, error) {
	dataSource := "user=" + env.GetDBUser() +
		" password=" + env.GetDBPassword() +
		" dbname=" + env.GetDBName() +
		" host=" + env.GetDBHost() +
		" port=" + env.GetDBPort() +
		" sslmode=disable"

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if env.GetDBMaxOpenConns() > 0 {
		db.SetMaxOpenConns(env.GetDBMaxOpenConns())
	}

	if env.GetDBMaxIdleConns() > 0 {
		db.SetMaxIdleConns(env.GetDBMaxIdleConns())
	}

	if env.GetDBMaxConnLifetime() > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(env.GetDBMaxConnLifetime()))
	}

	if env.GetDBMaxIdleConnLifetime() > 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(env.GetDBMaxIdleConnLifetime()))
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
