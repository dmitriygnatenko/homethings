package env

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
)

// nolint:gosec
const (
	envPath = "../../config/.env"

	appPortEnv = "APP_PORT"

	dbHostEnv     = "DB_HOST"
	dbPortEnv     = "DB_PORT"
	dbNameEnv     = "DB_NAME"
	dbUserEnv     = "DB_USER"
	dbPasswordEnv = "DB_PASSWORD"
)

type env struct {
	appPort string

	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
}

func Init() (interfaces.IEnv, error) {
	res := &env{}

	path := flag.String("config", envPath, "path to .env config")
	flag.Parse()

	file, err := os.Open(*path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}

		set(res, parts[0], parts[1])
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return res, err
}

func set(res *env, key string, value string) {
	switch key {
	case appPortEnv:
		res.appPort = value

	case dbHostEnv:
		res.dbHost = value
	case dbPortEnv:
		res.dbPort = value
	case dbNameEnv:
		res.dbName = value
	case dbUserEnv:
		res.dbUser = value
	case dbPasswordEnv:
		res.dbPassword = value
	}
}

func (e *env) GetAppPort() string {
	return e.appPort
}

func (e *env) GetDBHost() string {
	return e.dbHost
}

func (e *env) GetDBPort() string {
	return e.dbPort
}

func (e *env) GetDBName() string {
	return e.dbName
}

func (e *env) GetDBUser() string {
	return e.dbUser
}

func (e *env) GetDBPassword() string {
	return e.dbPassword
}
