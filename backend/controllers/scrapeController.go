package authController

import (
	"github.com/gofiber/fiber/v2"
	scrapeService "github.com/martynasd123/golang-scraper/services"
)

type ScrapeController struct {
	service *scrapeService.ScrapeService
}

func NewScrapeController(service *scrapeService.ScrapeService) *ScrapeController {
	return &ScrapeController{service: service}
}

func (*ScrapeController) AddTask(ctx *fiber.Ctx) error {
	return ctx.SendString("about")
}
