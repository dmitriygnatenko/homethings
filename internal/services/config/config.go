package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Service struct {
	appPort               uint16
	dbDriver              string
	dbHost                string
	dbPort                uint16
	dbName                string
	dbUser                string
	dbPassword            string
	dbMaxOpenConns        uint16
	dbMaxIdleConns        uint16
	dbMaxOpenConnLifetime time.Duration
	dbMaxIdleConnLifetime time.Duration
	corsAllowOrigins      string
	corsAllowMethods      string
	jwtSecretKey          string
	jwtLifeTime           time.Duration
	basicAuthUser         string
	basicAuthPassword     string
	smtpHost              string
	smtpPort              uint16
	smtpUser              string
	smtpPassword          string
	loggerStdoutEnabled   bool
	loggerStdoutLevel     string
	loggerEmailEnabled    bool
	loggerEmailLevel      string
	loggerEmailRecipient  string
	loggerEmailSubject    string
}

func Init() (*Service, error) {
	s := struct {
		AppPort               uint16        `envconfig:"APP_PORT"`
		DBDriver              string        `envconfig:"DB_DRIVER"`
		DBHost                string        `envconfig:"DB_HOST"`
		DBPort                uint16        `envconfig:"DB_PORT"`
		DBName                string        `envconfig:"DB_NAME"`
		DBUser                string        `envconfig:"DB_USER"`
		DBPassword            string        `envconfig:"DB_PASSWORD"`
		DBMaxOpenConns        uint16        `envconfig:"DB_MAX_OPEN_CONNS"`
		DBMaxIdleConns        uint16        `envconfig:"DB_MAX_IDLE_CONNS"`
		DBMaxOpenConnLifetime time.Duration `envconfig:"DB_MAX_OPEN_CONN_LIFETIME"`
		DBMaxIdleConnLifetime time.Duration `envconfig:"DB_MAX_IDLE_CONN_LIFETIME"`
		CORSAllowOrigins      string        `envconfig:"CORS_ALLOW_ORIGING"`
		CORSAllowMethods      string        `envconfig:"CORS_ALLOW_METHODS"`
		JWTSecretKey          string        `envconfig:"JWT_SECRET_KEY"`
		JWTLifeTime           time.Duration `envconfig:"JWT_LIFETIME"`
		BasicAuthUser         string        `envconfig:"BASIC_AUTH_USER"`
		BasicAuthPassword     string        `envconfig:"BASIC_AUTH_PASSWORD"`
		SMTPHost              string        `envconfig:"SMTP_HOST"`
		SMTPPort              uint16        `envconfig:"SMTP_PORT"`
		SMTPUser              string        `envconfig:"SMTP_USER"`
		SMTPPassword          string        `envconfig:"SMTP_PASSWORD"`
		LoggerStdoutEnabled   bool          `envconfig:"LOGGER_STDOUT_ENABLED"`
		LoggerStdoutLevel     string        `envconfig:"LOGGER_STDOUT_LEVEL"`
		LoggerEmailEnabled    bool          `envconfig:"LOGGER_EMAIL_ENABLED"`
		LoggerEmailLevel      string        `envconfig:"LOGGER_EMAIL_LEVEL"`
		LoggerEmailRecipient  string        `envconfig:"LOGGER_EMAIL_RECIPIENT"`
		LoggerEmailSubject    string        `envconfig:"LOGGER_EMAIL_SUBJECT"`
	}{}

	_ = godotenv.Overload()

	if err := envconfig.Process("", &s); err != nil {
		return nil, err
	}

	return &Service{
		appPort:               s.AppPort,
		dbDriver:              s.DBDriver,
		dbHost:                s.DBHost,
		dbPort:                s.DBPort,
		dbName:                s.DBName,
		dbUser:                s.DBUser,
		dbPassword:            s.DBPassword,
		dbMaxOpenConns:        s.DBMaxOpenConns,
		dbMaxIdleConns:        s.DBMaxIdleConns,
		dbMaxOpenConnLifetime: s.DBMaxOpenConnLifetime,
		dbMaxIdleConnLifetime: s.DBMaxIdleConnLifetime,
		corsAllowOrigins:      s.CORSAllowOrigins,
		corsAllowMethods:      s.CORSAllowMethods,
		jwtSecretKey:          s.JWTSecretKey,
		jwtLifeTime:           s.JWTLifeTime,
		basicAuthUser:         s.BasicAuthUser,
		basicAuthPassword:     s.BasicAuthPassword,
		smtpHost:              s.SMTPHost,
		smtpPort:              s.SMTPPort,
		smtpUser:              s.SMTPUser,
		smtpPassword:          s.SMTPPassword,
		loggerStdoutEnabled:   s.LoggerStdoutEnabled,
		loggerStdoutLevel:     s.LoggerStdoutLevel,
		loggerEmailEnabled:    s.LoggerEmailEnabled,
		loggerEmailLevel:      s.LoggerEmailLevel,
		loggerEmailRecipient:  s.LoggerEmailRecipient,
		loggerEmailSubject:    s.LoggerEmailSubject,
	}, nil
}

func (e *Service) AppPort() uint16 {
	return e.appPort
}

func (e *Service) DBDriver() string {
	return e.dbDriver
}

func (e *Service) DBHost() string {
	return e.dbHost
}

func (e *Service) DBPort() uint16 {
	return e.dbPort
}

func (e *Service) DBName() string {
	return e.dbName
}

func (e *Service) DBUser() string {
	return e.dbUser
}

func (e *Service) DBPassword() string {
	return e.dbPassword
}

func (e *Service) CORSAllowOrigins() string {
	return e.corsAllowOrigins
}

func (e *Service) CORSAllowMethods() string {
	return e.corsAllowMethods
}

func (e *Service) DBMaxOpenConns() uint16 {
	return e.dbMaxOpenConns
}

func (e *Service) DBMaxIdleConns() uint16 {
	return e.dbMaxIdleConns
}

func (e *Service) DBMaxOpenConnLifetime() time.Duration {
	return e.dbMaxOpenConnLifetime
}

func (e *Service) DBMaxIdleConnLifetime() time.Duration {
	return e.dbMaxIdleConnLifetime
}

func (e *Service) SMTPHost() string {
	return e.smtpHost
}

func (e *Service) SMTPPort() uint16 {
	return e.smtpPort
}

func (e *Service) SMTPUser() string {
	return e.smtpUser
}

func (e *Service) SMTPPassword() string {
	return e.smtpPassword
}

func (e *Service) JWTSecretKey() string {
	return e.jwtSecretKey
}

func (e *Service) JWTLifetime() time.Duration {
	return e.jwtLifeTime
}

func (e *Service) BasicAuthUser() string {
	return e.basicAuthUser
}

func (e *Service) BasicAuthPassword() string {
	return e.basicAuthPassword
}

func (e *Service) LoggerStdoutEnabled() bool {
	return e.loggerStdoutEnabled
}

func (e *Service) LoggerStdoutLevel() string {
	return e.loggerStdoutLevel
}

func (e *Service) LoggerEmailEnabled() bool {
	return e.loggerEmailEnabled
}

func (e *Service) LoggerEmailLevel() string {
	return e.loggerEmailLevel
}

func (e *Service) LoggerEmailRecipient() string {
	return e.loggerEmailRecipient
}

func (e *Service) LoggerEmailSubject() string {
	return e.loggerEmailSubject
}
