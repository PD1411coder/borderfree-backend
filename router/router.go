package router

import (
	"borderfree/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupSecuredRoutes(app *fiber.App) {
	api := app.Group("/api")

	products := api.Group("/products")
	products.Get("/", handler.GetAllProducts)
	products.Delete("/", handler.DeleteProduct)
	products.Post("/", handler.AddProduct)
	products.Put("/", handler.UpdateProduct)

	users := api.Group("/users")
	users.Get("/", handler.GetAllUsers)
}
