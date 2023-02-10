package interfaces

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

	GetCORSAllowOrigins() string
	GetCORSAllowMethods() string

	GetSMTPHost() string
	GetSMTPPort() string
	GetSMTPUser() string
	GetSMTPPassword() string

	GetJWTSecretKey() string
	GetJWTLifetime() int

	GetBasicAuthUser() string
	GetBasicAuthPassword() string

	GetErrorsEmail() string
}
