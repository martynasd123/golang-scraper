package scrape

type AddTaskResponse struct {
	Id int `json:"id"`
}

func CreateAddTaskResponse(id int) *AddTaskResponse {
	return &AddTaskResponse{id}
}
