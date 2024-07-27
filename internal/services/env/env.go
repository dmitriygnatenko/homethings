package env

import (
	"flag"

	"github.com/spf13/viper"
)

const defaultConfigPath = "../../.env"

type Service struct {
	appPort               string
	dbHost                string
	dbPort                string
	dbName                string
	dbUser                string
	dbPassword            string
	dbMaxOpenConns        int
	dbMaxIdleConns        int
	dbMaxConnLifetime     int
	dbMaxIdleConnLifetime int
	corsAllowOrigins      string
	corsAllowMethods      string
	jwtSecretKey          string
	jwtLifeTime           int
	basicAuthUser         string
	basicAuthPassword     string
	smtpHost              string
	smtpPort              string
	smtpUser              string
	smtpPassword          string
	errorsEmail           string
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
		AppPort               string `mapstructure:"APP_PORT"`
		DBHost                string `mapstructure:"DB_HOST"`
		DBPort                string `mapstructure:"DB_PORT"`
		DBName                string `mapstructure:"DB_NAME"`
		DBUser                string `mapstructure:"DB_USER"`
		DBPassword            string `mapstructure:"DB_PASSWORD"`
		DBMaxOpenConns        int    `mapstructure:"DB_MAX_OPEN_CONNS"`
		DBMaxIdleConns        int    `mapstructure:"DB_MAX_IDLE_CONNS"`
		DBMaxConnLifetime     int    `mapstructure:"DB_MAX_CONN_LIFETIME"`
		DBMaxIdleConnLifetime int    `mapstructure:"DB_MAX_IDLE_CONN_LIFETIME"`
		CORSAllowOrigins      string `mapstructure:"CORS_ALLOW_ORIGING"`
		CORSAllowMethods      string `mapstructure:"CORS_ALLOW_METHODS"`
		JWTSecretKey          string `mapstructure:"JWT_SECRET_KEY"`
		JWTLifeTime           int    `mapstructure:"JWT_LIFETIME"`
		BasicAuthUser         string `mapstructure:"BASIC_AUTH_USER"`
		BasicAuthPassword     string `mapstructure:"BASIC_AUTH_PASSWORD"`
		SMTPHost              string `mapstructure:"SMTP_HOST"`
		SMTPPort              string `mapstructure:"SMTP_PORT"`
		SMTPUser              string `mapstructure:"SMTP_USER"`
		SMTPPassword          string `mapstructure:"SMTP_PASSWORD"`
		ErrorsEmail           string `mapstructure:"ERRORS_EMAIL"`
	}{}

	if err := viper.Unmarshal(&s); err != nil {
		return nil, err
	}

	return &Service{
		appPort:               s.AppPort,
		dbHost:                s.DBHost,
		dbPort:                s.DBPort,
		dbName:                s.DBName,
		dbUser:                s.DBUser,
		dbPassword:            s.DBPassword,
		dbMaxOpenConns:        s.DBMaxOpenConns,
		dbMaxIdleConns:        s.DBMaxIdleConns,
		dbMaxConnLifetime:     s.DBMaxConnLifetime,
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
		errorsEmail:           s.ErrorsEmail,
	}, nil
}

func (e *Service) GetAppPort() string {
	return e.appPort
}

func (e *Service) GetDBHost() string {
	return e.dbHost
}

func (e *Service) GetDBPort() string {
	return e.dbPort
}

func (e *Service) GetDBName() string {
	return e.dbName
}

func (e *Service) GetDBUser() string {
	return e.dbUser
}

func (e *Service) GetDBPassword() string {
	return e.dbPassword
}

func (e *Service) GetCORSAllowOrigins() string {
	return e.corsAllowOrigins
}

func (e *Service) GetCORSAllowMethods() string {
	return e.corsAllowMethods
}

func (e *Service) GetDBMaxOpenConns() int {
	return e.dbMaxOpenConns
}

func (e *Service) GetDBMaxIdleConns() int {
	return e.dbMaxIdleConns
}

func (e *Service) GetDBMaxConnLifetime() int {
	return e.dbMaxConnLifetime
}

func (e *Service) GetDBMaxIdleConnLifetime() int {
	return e.dbMaxIdleConnLifetime
}

func (e *Service) GetSMTPHost() string {
	return e.smtpHost
}

func (e *Service) GetSMTPPort() string {
	return e.smtpPort
}

func (e *Service) GetSMTPUser() string {
	return e.smtpUser
}

func (e *Service) GetSMTPPassword() string {
	return e.smtpPassword
}

func (e *Service) GetJWTSecretKey() string {
	return e.jwtSecretKey
}

func (e *Service) GetJWTLifetime() int {
	return e.jwtLifeTime
}

func (e *Service) GetErrorsEmail() string {
	return e.errorsEmail
}

func (e *Service) GetBasicAuthUser() string {
	return e.basicAuthUser
}

func (e *Service) GetBasicAuthPassword() string {
	return e.basicAuthPassword
}
