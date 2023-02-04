package env

import (
	"git.dmitriygnatenko.ru/dima/homethings/internal/interfaces"
	"github.com/spf13/viper"
)

type env struct {
	AppPort string `mapstructure:"APP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	DBMaxOpenConns        int `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns        int `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxConnLifetime     int `mapstructure:"DB_MAX_CONN_LIFETIME"`
	DBMaxIdleConnLifetime int `mapstructure:"DB_MAX_IDLE_CONN_LIFETIME"`

	CORSAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGING"`
	CORSAllowMethods string `mapstructure:"CORS_ALLOW_METHODS"`

	JWTSecretKey string `mapstructure:"JWT_SECRET_KEY"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     string `mapstructure:"SMTP_PORT"`
	SMTPUser     string `mapstructure:"SMTP_USER"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	ErrorsEmail string `mapstructure:"ERRORS_EMAIL"`
}

func Init(configPath string) (interfaces.IEnv, error) {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	res := &env{}
	err := viper.Unmarshal(&res)

	return res, err
}

func (e *env) GetAppPort() string {
	return e.AppPort
}

func (e *env) GetDBHost() string {
	return e.DBHost
}

func (e *env) GetDBPort() string {
	return e.DBPort
}

func (e *env) GetDBName() string {
	return e.DBName
}

func (e *env) GetDBUser() string {
	return e.DBUser
}

func (e *env) GetDBPassword() string {
	return e.DBPassword
}

func (e *env) GetCORSAllowOrigins() string {
	return e.CORSAllowOrigins
}

func (e *env) GetCORSAllowMethods() string {
	return e.CORSAllowMethods
}

func (e *env) GetDBMaxOpenConns() int {
	return e.DBMaxOpenConns
}

func (e *env) GetDBMaxIdleConns() int {
	return e.DBMaxIdleConns
}

func (e *env) GetDBMaxConnLifetime() int {
	return e.DBMaxConnLifetime
}

func (e *env) GetDBMaxIdleConnLifetime() int {
	return e.DBMaxIdleConnLifetime
}

func (e *env) GetSMTPHost() string {
	return e.SMTPHost
}

func (e *env) GetSMTPPort() string {
	return e.SMTPPort
}

func (e *env) GetSMTPUser() string {
	return e.SMTPUser
}

func (e *env) GetSMTPPassword() string {
	return e.SMTPPassword
}

func (e *env) GetJWTSecretKey() string {
	return e.JWTSecretKey
}

func (e *env) GetErrorsEmail() string {
	return e.ErrorsEmail
}
