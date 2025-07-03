package sp

import (
	"git.dmitriygnatenko.ru/dima/go-common/db"
	"git.dmitriygnatenko.ru/dima/go-common/logger"
	"git.dmitriygnatenko.ru/dima/go-common/smtp"

	"git.dmitriygnatenko.ru/dima/homethings/internal/repositories"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/auth"
	"git.dmitriygnatenko.ru/dima/homethings/internal/services/config"
)

type ServiceProvider struct {
	config                      *config.Service
	auth                        *auth.Service
	tm                          *db.TxManager
	placeRepository             *repositories.PlaceRepository
	thingRepository             *repositories.ThingRepository
	tagRepository               *repositories.TagRepository
	placeThingRepository        *repositories.PlaceThingRepository
	placeImageRepository        *repositories.PlaceImageRepository
	thingImageRepository        *repositories.ThingImageRepository
	thingTagRepository          *repositories.ThingTagRepository
	userRepository              *repositories.UserRepository
	thingNotificationRepository *repositories.ThingNotificationRepository
	fileRepository              *repositories.FileRepository
}

func Init() (*ServiceProvider, error) {
	sp := &ServiceProvider{}

	// Init config
	configService, err := config.Init()
	if err != nil {
		return nil, err
	}
	sp.config = configService

	// Init auth
	authService, err := auth.Init(configService)
	if err != nil {
		return nil, err
	}
	sp.auth = authService

	// Init DB
	dbService, err := db.NewDB(
		db.NewConfig(
			db.WithDriver(configService.DBDriver()),
			db.WithHost(configService.DBHost()),
			db.WithPort(configService.DBPort()),
			db.WithUsername(configService.DBUser()),
			db.WithPassword(configService.DBPassword()),
			db.WithDatabase(configService.DBName()),
			db.WithMaxIdleConns(configService.DBMaxIdleConns()),
			db.WithMaxOpenConns(configService.DBMaxOpenConns()),
			db.WithMaxIdleConnLifetime(configService.DBMaxIdleConnLifetime()),
			db.WithMaxOpenConnLifetime(configService.DBMaxOpenConnLifetime()),
		),
	)
	if err != nil {
		return nil, err
	}

	// Init transaction manager
	sp.tm = db.NewTransactionManager(dbService)

	// Init repositories
	sp.placeRepository = repositories.InitPlaceRepository(dbService)
	sp.thingRepository = repositories.InitThingRepository(dbService)
	sp.placeThingRepository = repositories.InitPlaceThingRepository(dbService)
	sp.placeImageRepository = repositories.InitPlaceImageRepository(dbService)
	sp.thingImageRepository = repositories.InitThingImageRepository(dbService)
	sp.tagRepository = repositories.InitTagRepository(dbService)
	sp.thingTagRepository = repositories.InitThingTagRepository(dbService)
	sp.userRepository = repositories.InitUserRepository(dbService)
	sp.thingNotificationRepository = repositories.InitThingNotificationRepository(dbService)
	sp.fileRepository = repositories.InitFileRepository()

	// Init logger
	loggerConfigOptions := logger.ConfigOptions{
		logger.WithEmailLogEnabled(configService.LoggerEmailEnabled()),
		logger.WithEmailLogLevel(configService.LoggerEmailLevel()),
		logger.WithEmailSubject(configService.LoggerEmailSubject()),
		logger.WithEmailRecipient(configService.LoggerEmailRecipient()),
		logger.WithStdoutLogEnabled(configService.LoggerStdoutEnabled()),
		logger.WithStdoutLogLevel(configService.LoggerStdoutLevel()),
	}

	if len(configService.SMTPPassword()) > 0 &&
		len(configService.SMTPUser()) > 0 &&
		len(configService.SMTPHost()) > 0 &&
		configService.SMTPPort() > 0 {

		smtpService, smtpErr := smtp.NewSMTP(
			smtp.NewConfig(
				smtp.WithPassword(configService.SMTPPassword()),
				smtp.WithUsername(configService.SMTPUser()),
				smtp.WithHost(configService.SMTPHost()),
				smtp.WithPort(configService.SMTPPort()),
			),
		)
		if smtpErr != nil {
			return nil, smtpErr
		}

		loggerConfigOptions.Add(logger.WithSMTPClient(smtpService))
	}

	err = logger.Init(logger.NewConfig(loggerConfigOptions...))
	if err != nil {
		return nil, err
	}

	return sp, nil
}

func (sp *ServiceProvider) ConfigService() *config.Service {
	return sp.config
}

func (sp *ServiceProvider) AuthService() *auth.Service {
	return sp.auth
}

func (sp *ServiceProvider) TransactionManager() *db.TxManager {
	return sp.tm
}

func (sp *ServiceProvider) PlaceRepository() *repositories.PlaceRepository {
	return sp.placeRepository
}

func (sp *ServiceProvider) ThingRepository() *repositories.ThingRepository {
	return sp.thingRepository
}

func (sp *ServiceProvider) TagRepository() *repositories.TagRepository {
	return sp.tagRepository
}

func (sp *ServiceProvider) PlaceThingRepository() *repositories.PlaceThingRepository {
	return sp.placeThingRepository
}

func (sp *ServiceProvider) PlaceImageRepository() *repositories.PlaceImageRepository {
	return sp.placeImageRepository
}

func (sp *ServiceProvider) ThingImageRepository() *repositories.ThingImageRepository {
	return sp.thingImageRepository
}

func (sp *ServiceProvider) ThingTagRepository() *repositories.ThingTagRepository {
	return sp.thingTagRepository
}

func (sp *ServiceProvider) ThingNotificationRepository() *repositories.ThingNotificationRepository {
	return sp.thingNotificationRepository
}

func (sp *ServiceProvider) UserRepository() *repositories.UserRepository {
	return sp.userRepository
}

func (sp *ServiceProvider) FileRepository() *repositories.FileRepository {
	return sp.fileRepository
}
