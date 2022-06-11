package main

import (
	"borderfree/handler"
	"borderfree/router"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())

	auth := app.Group("/auth")
	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("ANiceLittleSecretOfMine987654321"),
	}))

	router.SetupSecuredRoutes(app)

	app.Listen(":" + os.Getenv("PORT"))
}
