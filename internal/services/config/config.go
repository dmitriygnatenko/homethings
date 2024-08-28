package config

import (
	"flag"
	"time"

	"github.com/spf13/viper"
)

const defaultConfigPath = "./.env"

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
	dbSSLMode             string
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
	loggerStdoutAddSource bool
	loggerEmailEnabled    bool
	loggerEmailLevel      string
	loggerEmailAddSource  bool
	loggerEmailRecipient  string
	loggerEmailSubject    string
}

func Init() (*Service, error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to .env config file")
	flag.Parse()

	if configPath == "" {
		configPath = defaultConfigPath
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	s := struct {
		AppPort               uint16        `mapstructure:"APP_PORT"`
		DBDriver              string        `mapstructure:"DB_DRIVER"`
		DBHost                string        `mapstructure:"DB_HOST"`
		DBPort                uint16        `mapstructure:"DB_PORT"`
		DBName                string        `mapstructure:"DB_NAME"`
		DBUser                string        `mapstructure:"DB_USER"`
		DBPassword            string        `mapstructure:"DB_PASSWORD"`
		DBSSLMode             string        `mapstructure:"DB_SSL_MODE"`
		DBMaxOpenConns        uint16        `mapstructure:"DB_MAX_OPEN_CONNS"`
		DBMaxIdleConns        uint16        `mapstructure:"DB_MAX_IDLE_CONNS"`
		DBMaxOpenConnLifetime time.Duration `mapstructure:"DB_MAX_OPEN_CONN_LIFETIME"`
		DBMaxIdleConnLifetime time.Duration `mapstructure:"DB_MAX_IDLE_CONN_LIFETIME"`
		CORSAllowOrigins      string        `mapstructure:"CORS_ALLOW_ORIGING"`
		CORSAllowMethods      string        `mapstructure:"CORS_ALLOW_METHODS"`
		JWTSecretKey          string        `mapstructure:"JWT_SECRET_KEY"`
		JWTLifeTime           time.Duration `mapstructure:"JWT_LIFETIME"`
		BasicAuthUser         string        `mapstructure:"BASIC_AUTH_USER"`
		BasicAuthPassword     string        `mapstructure:"BASIC_AUTH_PASSWORD"`
		SMTPHost              string        `mapstructure:"SMTP_HOST"`
		SMTPPort              uint16        `mapstructure:"SMTP_PORT"`
		SMTPUser              string        `mapstructure:"SMTP_USER"`
		SMTPPassword          string        `mapstructure:"SMTP_PASSWORD"`
		LoggerStdoutEnabled   bool          `mapstructure:"LOGGER_STDOUT_ENABLED"`
		LoggerStdoutLevel     string        `mapstructure:"LOGGER_STDOUT_LEVEL"`
		LoggerStdoutAddSource bool          `mapstructure:"LOGGER_STDOUT_ADD_SOURCE"`
		LoggerEmailEnabled    bool          `mapstructure:"LOGGER_EMAIL_ENABLED"`
		LoggerEmailLevel      string        `mapstructure:"LOGGER_EMAIL_LEVEL"`
		LoggerEmailAddSource  bool          `mapstructure:"LOGGER_EMAIL_ADD_SOURCE"`
		LoggerEmailRecipient  string        `mapstructure:"LOGGER_EMAIL_RECIPIENT"`
		LoggerEmailSubject    string        `mapstructure:"LOGGER_EMAIL_SUBJECT"`
	}{}

	if err := viper.Unmarshal(&s); err != nil {
		return nil, err
	}

	return &Service{
		appPort:               s.AppPort,
		dbDriver:              s.DBDriver,
		dbSSLMode:             s.DBSSLMode,
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
		loggerStdoutAddSource: s.LoggerStdoutAddSource,
		loggerEmailEnabled:    s.LoggerEmailEnabled,
		loggerEmailLevel:      s.LoggerEmailLevel,
		loggerEmailAddSource:  s.LoggerEmailAddSource,
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

func (e *Service) DBSSLMode() string {
	return e.dbSSLMode
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

func (e *Service) LoggerStdoutAddSource() bool {
	return e.loggerStdoutAddSource
}

func (e *Service) LoggerEmailEnabled() bool {
	return e.loggerEmailEnabled
}

func (e *Service) LoggerEmailLevel() string {
	return e.loggerEmailLevel
}

func (e *Service) LoggerEmailAddSource() bool {
	return e.loggerEmailAddSource
}

func (e *Service) LoggerEmailRecipient() string {
	return e.loggerEmailRecipient
}

func (e *Service) LoggerEmailSubject() string {
	return e.loggerEmailSubject
}
