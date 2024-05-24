package main

import (
	"github.com/gofiber/fiber/v2"
	controllers "github.com/martynasd123/golang-scraper/controllers"
	services "github.com/martynasd123/golang-scraper/services"
)

func DefineAuthRoutes(router fiber.Router) {
	service := services.NewAuthService()
	controller := controllers.NewAuthController(service)

	router.Post("/", controller.Authenticate)
	router.Post("/refresh-token", controller.RefreshToken)
	router.Post("/logout", controller.LogOut)
}

func DefineScrapeRoutes(router fiber.Router) {
	service := services.NewScrapeService()
	controller := controllers.NewScrapeController(service)

	router.Post("/add-task", controller.AddTask)
}

func DefineRoutes(router fiber.Router) {
	DefineAuthRoutes(router.Group("/auth"))
	DefineScrapeRoutes(router.Group("/scrape"))
}

func main() {
	app := fiber.New()

	DefineRoutes(app.Group("/api"))

	app.Listen(":3000")
}
