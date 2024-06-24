package main

import (
	"log"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func (srv *Service) init() {
	srv.MongoDB = srv.initMongo()
	srv.Router = srv.initRouter()

	appRouter := srv.Router.Group(srv.Config.SERVICE_BASEPATH)
	{
		appRouter.Post("/login", srv.Login)
		appRouter.Post("/register", srv.UserRegistration)
		appRouter.Get("/fieldCatalogues", srv.FieldCatalogues)
		appRouter.Post("/sendOTPVerificationMail", srv.SendOTPVerificationMail)
		appRouter.Post("/verifyOTP", srv.VerifyOTP)
		appRouter.Use(ValidateToken(srv))
		appRouter.Post("/create", srv.Create)
		appRouter.Put("/update/:id", srv.Update)
		appRouter.Put("/updatePassword/:id", srv.UpdatePassword)
		appRouter.Delete("/delete/:id/:hardDelete", srv.Delete)
		appRouter.Get("/get", srv.Get)
	}
}

// @Console - USER Service
// @version 1.0
// @description User Microservice for Console

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	userSvc := &Service{}
	var err error
	userSvc.Config, err = loadConfig()
	if err != nil {
		log.Fatal(err.Error(), err)
	}
	userSvc.init()
	userSvc.Router.Listen(":" + userSvc.Config.SERVICE_PORT)
}
