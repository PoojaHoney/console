package main

import (
	"context"
	"fmt"
	"log"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"net/http"
	docs "products/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func loadConfig() (config ServiceConfigurations, err error) {
	// Set the configuration type to YAML
	// viper.SetConfigType("yaml")
	// // Set the configuration path to /etc/config
	// viper.AddConfigPath("/etc/config")
	// // Set the configuration name to lms-qa-console-config
	// viper.SetConfigName("lms-qa-console-config")

	/* In local  */
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Load environment variables
	viper.AutomaticEnv()
	// Read the configuration file
	errOccured := viper.ReadInConfig()
	if errOccured != nil {
		return config, errOccured
	}
	// Unmarshal the configuration into the config struct
	errOccured = viper.Unmarshal(&config)
	if errOccured != nil {
		return config, errOccured
	}
	return
}

// initRouter initializes and configures the router for the service.
func (srv *Service) initRouter() *fiber.App {
	// Create a new Fiber app instance
	app := fiber.New(fiber.Config{
		ProxyHeader:             fiber.HeaderXForwardedFor,
		AppName:                 srv.Config.SERVICE_NAME,
		BodyLimit:               1 * 1024 * 1024 * 1024,
		EnableTrustedProxyCheck: true,
	})
	// Use the logger middleware
	app.Use(logger.New())
	// Use the recover middleware
	app.Use(recover.New())
	// Use the CORS middleware for handling cross-origin requests
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "POST, GET, PUT, DELETE",
		AllowHeaders:     "Access-Control-Allow-Origin, Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length",
		MaxAge:           86400,
	}))
	// Use the compress middleware for response compression
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	// Define a route for the health check endpoint
	app.Get(srv.Config.SERVICE_BASEPATH+"/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(&fiber.Map{"status": "Ok"})
	})
	// Define a route for the metrics endpoint
	app.Get(srv.Config.SERVICE_BASEPATH+"/metrics", monitor.New(monitor.Config{Title: "Config service Metrics Page"}))

	docs.SwaggerInfo.Title = "Swagger for Console Products Processor APIs"
	docs.SwaggerInfo.Description = "Console Products Processor API's"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", srv.Config.SERVICE_HOST, srv.Config.SERVICE_PORT)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	app.Get(fmt.Sprintf("%s/swagger/*", srv.Config.SERVICE_BASEPATH), swagger.HandlerDefault)
	return app
}

func (srv *Service) initMongo() *mongo.Client {
	var clientOptions *options.ClientOptions
	if srv.Config.MONGO_USER != "" && srv.Config.MONGO_PASSWORD != "" {
		clientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", srv.Config.MONGO_USER, srv.Config.MONGO_PASSWORD, srv.Config.MONGO_HOST, srv.Config.MONGO_PORT))
	} else {
		clientOptions = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", srv.Config.MONGO_HOST, srv.Config.MONGO_PORT))
	}
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("Unable to create database client.", err)
	}
	fmt.Println("DB Connected")
	return client
}
