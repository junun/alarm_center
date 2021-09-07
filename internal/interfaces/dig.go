package interfaces

import (
	"alarm_center/internal/application/service"
	"alarm_center/internal/config"
	"alarm_center/internal/infras/db"
	"alarm_center/internal/infras/persistence"
	"go.uber.org/dig"
)

var container = dig.New()

func BuildContainer() *dig.Container {
	// config
	container.Provide(config.LoadConfig)

	// init
	container.Provide(config.InitAppConf)

	// init email
	container.Provide(config.InitEmailConf)

	// DB
	container.Provide(config.InitDBConf)

	// cron
	container.Provide(db.NewCron)

	// sonar
	container.Provide(db.InitSonarClient)

	// redis
	container.Provide(config.InitRedisConn)

	// user
	container.Provide(persistence.NewUserRepository)
	container.Provide(service.NewUserService)

	// dingtalk
	container.Provide(persistence.NewDingtalkRepository)
	container.Provide(service.NewDingtalkService)

	//email
	container.Provide(persistence.NewEmailRepository)
	container.Provide(service.NewEmailService)

	//cron
	container.Provide(persistence.NewCronRepository)
	container.Provide(service.NewCronService)

	// sonar
	container.Provide(persistence.NewSonarRepository)
	container.Provide(service.NewSonarService)

	container.Provide(NewServer)

	return container
}

func Invoke(i interface{}) error {
	return container.Invoke(i)
}