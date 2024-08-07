package main

import (
	"fmt"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"net/http"
	docs "users/docs"
	// "time"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/gofiber/fiber/v2/middleware/session"
	// "github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// loadConfig loads the service configuration from a YAML file and environment variables.
// It returns the loaded configuration or an error if the loading process fails.
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
	// Initialize session middleware
	// store := session.New(session.Config{
	// 	KeyLookup:      "cookie:session_id",
	// 	CookieSecure:   true,
	// 	CookieHTTPOnly: true,
	// 	CookieSameSite: "Strict",
	// 	Expiration:     24 * time.Hour,
	// })
	// // Middleware to handle sessions
	// app.Use(func(c *fiber.Ctx) error {
	// 	sess, err := store.Get(c)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	defer sess.Save()
	// 	c.Locals("session", sess)
	// 	return c.Next()
	// })

	// // Optional: Initialize CSRF protection middleware
	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-CSRF-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Strict",
	// 	CookieHTTPOnly: true,
	// 	}))
	// Define a route for the health check endpoint
	app.Get(srv.Config.SERVICE_BASEPATH+"/health", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON(&fiber.Map{"status": "Ok"})
	})
	// Define a route for the metrics endpoint
	app.Get(srv.Config.SERVICE_BASEPATH+"/metrics", monitor.New(monitor.Config{Title: "Config service Metrics Page"}))

	docs.SwaggerInfo.Title = "Swagger for Console User Processor APIs"
	docs.SwaggerInfo.Description = "Console User Processor API's"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", srv.Config.SERVICE_HOST, srv.Config.SERVICE_PORT)
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	app.Get(fmt.Sprintf("%s/swagger/*", srv.Config.SERVICE_BASEPATH), swagger.HandlerDefault)
	return app
}

func (srv *Service) initPostgres() *gorm.DB {
	dns := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		srv.Config.POSTGRES_HOST, srv.Config.POSTGRES_USERNAME, srv.Config.POSTGRES_PASSWORD, srv.Config.POSTGRES_DATABASE, srv.Config.POSTGRES_PORT)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Enable UUID generation
    db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&User{}, &Address{}, &UserOTP{}, &Tokens{})
	return db
}
