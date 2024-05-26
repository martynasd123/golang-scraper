package scrape

import (
	. "github.com/martynasd123/golang-scraper/storage"
)

type AddTaskResponse struct {
	Id int `json:"id"`
}

func CreateAddTaskResponse(id int) *AddTaskResponse {
	return &AddTaskResponse{id}
}

type TaskStatusResponse struct {
	Id                *int    `json:"id"`
	Status            string  `json:"status"`
	Link              string  `json:"link"`
	ExternalLinks     *int    `json:"externalLinks"`
	InternalLinks     *int    `json:"internalLinks"`
	InaccessibleLinks *int    `json:"inaccessibleLinks"`
	HtmlVersion       *string `json:"htmlVersion"`
	PageTitle         *string `json:"pageTitle"`
	HeadingsByLevel   *[6]int `json:"headingsByLevel"`
	CrawledLinks      int     `json:"crawledLinks"`
	Error             *string `json:"error"`
}

func CreateTaskStatusResponse(task *Task) *TaskStatusResponse {
	response := &TaskStatusResponse{}
	response.Id = task.Id
	response.Status = task.Status
	response.InaccessibleLinks = task.InaccessibleLinks
	response.HtmlVersion = task.HtmlVersion
	response.PageTitle = task.PageTitle
	response.HeadingsByLevel = task.HeadingsByLevel
	response.Link = task.Link.String()
	response.ExternalLinks = task.ExternalLinks
	response.InternalLinks = task.InternalLinks
	response.CrawledLinks = task.CrawledLinks
	response.Error = task.Error
	return response
}

type TaskListItem struct {
	Id                *int    `json:"id"`
	Link              string  `json:"link"`
	InaccessibleLinks *int    `json:"inaccessibleLinks"`
	PageTitle         *string `json:"pageTitle"`
	CrawledLinks      int     `json:"crawledLinks"`
	Status            string  `json:"status"`
	Error             *string `json:"error"`
}

func CreateTaskListItemResponse(task *Task) *TaskListItem {
	response := &TaskListItem{}
	response.Id = task.Id
	response.InaccessibleLinks = task.InaccessibleLinks
	response.PageTitle = task.PageTitle
	response.Link = task.Link.String()
	response.CrawledLinks = task.CrawledLinks
	response.Error = task.Error
	response.Status = task.Status
	return response
}
