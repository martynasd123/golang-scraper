package authController

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	request "github.com/martynasd123/golang-scraper/models/request"
	response "github.com/martynasd123/golang-scraper/models/response"
	scrapeService "github.com/martynasd123/golang-scraper/services/scrape"
)

type ScrapeController struct {
	service *scrapeService.ScrapeService
}

func CreateScrapeController(service *scrapeService.ScrapeService) *ScrapeController {
	return &ScrapeController{service: service}
}

func (controller *ScrapeController) AddTask(ctx *gin.Context) {
	var body request.AddTaskRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		// todo
		return
	}

	parsedUrl, err := url.Parse(body.Link)
	if err != nil {
		// todo
		return
	}

	if parsedUrl.Scheme != "https" && parsedUrl.Scheme != "http" {
		// Only http/https is supported
		// todo
		return
	}

	if parsedUrl.Fragment != "" {
		// Ignore fragments
		parsedUrl.Fragment = ""
	}

	id, err := controller.service.AddTask(parsedUrl)
	if err != nil {
		log.Printf("error occurred when adding task: %v", err)
	}

	ctx.JSON(http.StatusOK, response.CreateAddTaskResponse(id))
}

func (controller *ScrapeController) Listen(ctx *gin.Context) {
	taskId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(400, "invalid task id")
		return
	}
	done := ctx.Writer.CloseNotify()

	err, data, notifyDone := controller.service.RegisterListener(taskId)

	if err != nil {
		if err.Error() == "task already finished" {
			task, _ := controller.service.GetTaskById(taskId)
			ctx.SSEvent("update", response.CreateTaskStatusResponse(task))
		} else if strings.HasPrefix(err.Error(), "no task found with id") {
			ctx.String(400, "invalid task id")
		} else {
			log.Printf("unexpected error occurred while attempting to register listener for task: %s", err)
			ctx.String(500, "unexpected error occurred")
		}
		return
	}

	// Stream response
	ctx.Stream(func(w io.Writer) bool {
		for {
			select {
			case <-done:
				// Client closed connection
				notifyDone <- struct{}{}
				return false // False to indicate no more data should be sent
			case task, ok := <-data:
				if !ok {
					// No more data
					return false
				}
				ctx.SSEvent("update", response.CreateTaskStatusResponse(&task))
				return true
			}
		}
	})
}
