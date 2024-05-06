package main

// import "github.com/gofiber/fiber/v2"/

type ServiceConfigurations struct {
	SERVICE_PORT     string `mapstructure:"SERVICE_PORT"`
	SERVICE_NAME     string `mapstructure:"SERVICE_NAME"`
	SERVICE_DOMAIN   string `mapstructure:"SERVICE_DOMAIN"`
	SERVICE_BASEPATH string `mapstructure:"SERVICE_BASEPATH"`
	SERVICE_VERSION  string `mapstructure:"SERVICE_VERSION"`
	TOKEN_COLLECTION string `mapstructure:"TOKEN_COLLECTION"`
	MONGO_DATABASE   string `mapstructure:"MONGO_DATABASE"`
	MONGO_HOST       string `mapstructure:"MONGO_HOST"`
	MONGO_PORT       string `mapstructure:"MONGO_PORT"`
	MONGO_USER       string `mapstructure:"MONGO_USER"`
	MONGO_PASSWORD   string `mapstructure:"MONGO_PASSWORD"`
	JAEGER_URL       string `mapstructure:"JAEGER_URL"`
	TRACER_ENABLED   bool   `mapstructure:"TRACER_ENABLED"`
	SMTP_SERVER      string `mapstructure:"SMTP_SERVER"`
	SMTP_BASEEMAIL   string `mapstructure:"SMTP_BASEEMAIL"`
	SMTP_PASSKEY     string `mapstructure:"SMTP_PASSKEY"`
	SMTP_PORT        string `mapstructure:"SMTP_PORT"`
}

// type Service struct {
// 	Router *fiber.App
// 	Config
// }
