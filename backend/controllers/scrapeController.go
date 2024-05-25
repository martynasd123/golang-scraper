package authController

import (
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	request "github.com/martynasd123/golang-scraper/models/request"
	response "github.com/martynasd123/golang-scraper/models/response"
	scrapeService "github.com/martynasd123/golang-scraper/services/scrape"
)

type ScrapeController struct {
	service *scrapeService.ScrapeService
}

func NewScrapeController(service *scrapeService.ScrapeService) *ScrapeController {
	return &ScrapeController{service: service}
}

func (controller *ScrapeController) AddTask(ctx *gin.Context) {
	var body request.AddTaskRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		// todo
		return
	}

	url, err := url.Parse(body.Link)
	if err != nil {
		// todo
		return
	}

	if url.Scheme != "https" && url.Scheme != "http" {
		// Only http/https is supported
		// todo
		return
	}

	if url.Fragment != "" {
		// Ignore fragments
		url.Fragment = ""
	}

	id, err := controller.service.AddTask(url)
	if err != nil {
		// todo handle error
	}

	ctx.JSON(http.StatusOK, response.CreateAddTaskResponse(id))
}

func (controller *ScrapeController) Listen(ctx *gin.Context) {

	taskId := 0
	done := ctx.Writer.CloseNotify()
	listenerId, channel, _ := controller.service.RegisterListener(taskId)

	ctx.Stream(func(w io.Writer) bool {
		for {
			select {
			case <-done:
				controller.service.UnregisterListener(taskId, listenerId)
				return true
			case <-channel:
				ctx.SSEvent("update", "value")
				return true
			}
		}
	})
}
