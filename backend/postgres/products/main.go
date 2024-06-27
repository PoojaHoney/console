package main

import (
	"fmt"

	"github.com/go-playground/validator"
)

var validate = validator.New()

func (srv *Service) init() {
	srv.MongoDB = srv.initMongo()
	srv.Router = srv.initRouter()

	appRouter := srv.Router.Group(srv.Config.SERVICE_BASEPATH)
	{
		// appRouter.Use(ValidateToken(srv))
		appRouter.Post("/createProduct", srv.CreateProduct)
		appRouter.Post("/createConfiguration", srv.CreateConfiguration)
		appRouter.Post("/createResource", srv.CreateResource)
		appRouter.Post("/createPlan", srv.CreatePlan)
		appRouter.Post("/createMicroService", srv.CreateMicroService)
		appRouter.Post("/createVersion", srv.CreateVersion)
		appRouter.Get("/getFullProductDetails", srv.GetFullProductDetails)
		appRouter.Post("/activateProduct", srv.ActivateProduct)
	}
}

// @Console - PRODUCT Service
// @version 1.0
// @description Product Microservice for Console
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication
func main() {
	productSrv := &Service{}
	var err error
	productSrv.Config, err = loadConfig()
	if err != nil {
		panic(err)
	}
	productSrv.MongoDB = productSrv.initMongo()
	productSrv.Router = productSrv.initRouter()
	productSrv.init()
	productSrv.Router.Listen(fmt.Sprintf(":%s", productSrv.Config.SERVICE_PORT))
}
