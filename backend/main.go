package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/martynasd123/golang-scraper/controllers"
	. "github.com/martynasd123/golang-scraper/services/auth"
	. "github.com/martynasd123/golang-scraper/services/scrape"
	"github.com/martynasd123/golang-scraper/storage"
	"log"
)

type ApplicationContext struct {
	AuthService   *AuthService
	ScrapeService *ScrapeService

	AuthController   *AuthController
	ScrapeController *ScrapeController

	AuthDao storage.AuthDao
	TaskDao storage.TaskDao

	RequireAuthMiddleware gin.HandlerFunc
}

func WireContext() *ApplicationContext {
	ctx := new(ApplicationContext)

	ctx.AuthDao = storage.CreateAuthInMemoryDao()
	ctx.TaskDao = storage.CreateTaskInMemoryDao()

	ctx.AuthService = CreateAuthService(ctx.AuthDao)
	ctx.ScrapeService = CreateTaskService(ctx.TaskDao)

	ctx.AuthController = CreateAuthController(ctx.AuthService)
	ctx.ScrapeController = CreateScrapeController(ctx.ScrapeService)

	ctx.RequireAuthMiddleware = RequireAuth(ctx.AuthService)
	return ctx
}

func DefineAuthRoutes(router *gin.RouterGroup, context *ApplicationContext) {
	router.POST("/", context.AuthController.Authenticate)
	router.POST("/refresh-token", context.AuthController.RefreshToken)
	router.POST("/log-out", context.RequireAuthMiddleware, context.AuthController.LogOut)
}

func DefineScrapeRoutes(router *gin.RouterGroup, context *ApplicationContext) {
	router.POST("/add-task", context.ScrapeController.AddTask)
	router.POST("/task/:id/interrupt", context.ScrapeController.InterruptTask)
	router.GET("/task/:id/listen", context.ScrapeController.Listen)
	router.GET("/tasks", context.ScrapeController.GetAllTasks)
}

func DefineRoutes(router *gin.RouterGroup, context *ApplicationContext) {
	DefineAuthRoutes(router.Group("/auth"), context)

	scrapeGroup := router.Group("/scrape")
	scrapeGroup.Use(context.RequireAuthMiddleware)
	DefineScrapeRoutes(scrapeGroup, context)
}

func main() {
	context := WireContext()

	// For testing purposes
	err := context.AuthService.CreateFakeUser("username", "password")
	if err != nil {
		log.Fatalln("Failed to create fake user:", err)
	}

	app := gin.Default()

	DefineRoutes(app.Group("/api"), context)

	app.Run()
}
