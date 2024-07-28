package db

import (
	"database/sql"
	"time"
)

type Env interface {
	AppPort() string
	DBHost() string
	DBPort() string
	DBName() string
	DBUser() string
	DBPassword() string
	DBMaxOpenConns() int
	DBMaxIdleConns() int
	DBMaxConnLifetime() int
	DBMaxIdleConnLifetime() int
}

func Init(env Env) (*sql.DB, error) {
	dataSource := "user=" + env.DBUser() +
		" password=" + env.DBPassword() +
		" dbname=" + env.DBName() +
		" host=" + env.DBHost() +
		" port=" + env.DBPort() +
		" sslmode=disable"

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		return nil, err
	}

	if env.DBMaxOpenConns() > 0 {
		db.SetMaxOpenConns(env.DBMaxOpenConns())
	}

	if env.DBMaxIdleConns() > 0 {
		db.SetMaxIdleConns(env.DBMaxIdleConns())
	}

	if env.DBMaxConnLifetime() > 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(env.DBMaxConnLifetime()))
	}

	if env.DBMaxIdleConnLifetime() > 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(env.DBMaxIdleConnLifetime()))
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
