package storage

import (
	"errors"
	"net/url"
	"sync"
)

type Task struct {
	Id                *int
	Link              url.URL
	Status            string
	ExternalLinks     *int
	InternalLinks     *int
	InaccessibleLinks *int
	HtmlVersion       *string
	PageTitle         *string
	HeadingsByLevel   *[6]int
	LoginFormPresent  *bool
	CrawledLinks      int
	Error             *string
}

func CreateTaskInitial(status string, link *url.URL) *Task {
	return &Task{
		Id:                nil,
		Link:              *link,
		Status:            status,
		ExternalLinks:     nil,
		InternalLinks:     nil,
		InaccessibleLinks: nil,
		HtmlVersion:       nil,
		PageTitle:         nil,
		HeadingsByLevel:   nil,
	}
}

// TaskDao as an interface to some storage mechanism (database, in-memory, some other...)
type TaskDao interface {

	// StoreTask stores tasks in the database. If the task is provided with an ID, overwrite existing task with
	// that id. If id is nil, insert a new task.
	// Returns:
	//   int: ID of the task
	StoreTask(task *Task) (int, error)

	// RetrieveTaskById retrieves task by given ID, or nil if not found
	RetrieveTaskById(id int) (*Task, error)
}

// TaskInMemoryDao is a simple in-memory storage mechanism for tasks.
// This likely shouldn't be used outside of testing environment,
// because it offers no persistence and is not very performant.
type TaskInMemoryDao struct {
	tasks  map[int]Task
	lastId int
	mu     sync.RWMutex
}

func (storage *TaskInMemoryDao) StoreTask(task *Task) (int, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if task.Id != nil {
		if _, ok := storage.tasks[*task.Id]; ok {
			// valid task id - overwrite
			storage.tasks[*task.Id] = *task
			return *task.Id, nil
		} else {
			return 0, errors.New("task with ID provided, but task does not exist")
		}
	}
	storage.lastId = storage.lastId + 1
	newId := storage.lastId
	storage.tasks[newId] = *task
	task.Id = &newId
	return newId, nil
}

func (storage *TaskInMemoryDao) RetrieveTaskById(id int) (*Task, error) {
	storage.mu.RLock()
	if value, ok := storage.tasks[id]; ok {
		return &value, nil
	}
	defer storage.mu.RUnlock()
	return nil, errors.New("no task found with id " + string(id))
}

func CreateTaskInMemoryDao() *TaskInMemoryDao {
	return &TaskInMemoryDao{tasks: make(map[int]Task), lastId: 0}
}
