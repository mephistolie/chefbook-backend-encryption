package config

import (
	"github.com/mephistolie/chefbook-backend-common/log"
	amqpConfig "github.com/mephistolie/chefbook-backend-common/mq/config"
)

const (
	EnvDev  = "develop"
	EnvProd = "production"
)

type Config struct {
	Environment *string
	Port        *int
	LogsPath    *string

	AuthService    AuthService
	ProfileService ProfileService
	RecipeService  RecipeService

	Database Database
	Amqp     amqpConfig.Amqp
	Smtp     Smtp
}

type AuthService struct {
	Addr *string
}

type ProfileService struct {
	Addr *string
}

type RecipeService struct {
	Addr *string
}

type Database struct {
	Host     *string
	Port     *int
	User     *string
	Password *string
	DBName   *string
}

type Smtp struct {
	Host         *string
	Port         *int
	Email        *string
	Password     *string
	SendAttempts *int
}

func (c Config) Validate() error {
	if *c.Environment != EnvProd {
		*c.Environment = EnvDev
	}
	return nil
}

func (c Config) Print() {
	log.Infof("ENCRYPTION SERVICE CONFIGURATION\n"+
		"Environment: %v\n"+
		"Port: %v\n"+
		"Logs path: %v\n\n"+
		"Auth Service Address: %v\n\n"+
		"Recipe Service Address: %v\n\n"+
		"Database host: %v\n"+
		"Database port: %v\n"+
		"Database name: %v\n\n"+
		"MQ host: %v\n"+
		"MQ port: %v\n"+
		"MQ vhost: %v\n\n"+
		"SMTP host: %v\n"+
		"SMTP port: %v\n\n",
		*c.Environment, *c.Port, *c.LogsPath,
		*c.AuthService.Addr, *c.RecipeService.Addr,
		*c.Database.Host, *c.Database.Port, *c.Database.DBName,
		*c.Amqp.Host, *c.Amqp.Port, *c.Amqp.VHost,
		*c.Smtp.Host, *c.Smtp.Port,
	)
}
