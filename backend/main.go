package main

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/martynasd123/golang-scraper/controllers"
	authService "github.com/martynasd123/golang-scraper/services/auth"
	scrapeService "github.com/martynasd123/golang-scraper/services/scrape"
	"github.com/martynasd123/golang-scraper/services/scrape/storage"
)

func DefineAuthRoutes(router *gin.RouterGroup) {
	service := authService.NewAuthService()
	controller := controllers.NewAuthController(service)

	router.POST("/", controller.Authenticate)
	router.POST("/refresh-token", controller.RefreshToken)
	router.POST("/logout", controller.LogOut)
}

func DefineScrapeRoutes(router *gin.RouterGroup) {
	inMemoryStorage := storage.CreateInMemoryStorage()
	service := scrapeService.NewScrapeService(inMemoryStorage)
	controller := controllers.NewScrapeController(service)

	router.POST("/add-task", controller.AddTask)
	router.GET("/:id/listen", controller.Listen)
}

func DefineRoutes(router *gin.RouterGroup) {
	DefineAuthRoutes(router.Group("/auth"))
	DefineScrapeRoutes(router.Group("/scrape"))
}

func main() {
	app := gin.Default()

	DefineRoutes(app.Group("/api"))

	app.Run()
}
