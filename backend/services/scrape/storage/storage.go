package storage

import (
	"errors"
	"net/url"
	"sync"
)

const (
	TASK_STATUS_QUEUED = iota
	TASK_STATUS_IN_PROGRESS
	TASK_STATUS_FINISHED
)

type Task struct {
	id     *int
	status int
	link   url.URL
	// ... rest of data
}

func CreateTask(status int, link *url.URL) *Task {
	return &Task{nil, status, *link}
}

// Interface to some storage mechanism (database, in-memory, some other...)
type TaskStorage interface {

	// Store task in the database. If the task is provided with an ID, overwrite existing task with
	// that id. If id is nil, insert a new task.
	// Returns:
	//   int: ID of the task
	StoreTask(task *Task) (int, error)

	// Retrieve task by given ID, or nil if not found
	RetrieveTaskById(id int) (*Task, error)
}

// A simple in-memory storage mechanism for tasks.
// This likely shouldn't be used outside of testing environment,
// because it offers no persistence and is not very performant.
type InMemoryStorage struct {
	tasks  map[int]Task
	lastId int
	mu     sync.RWMutex
}

func (storage *InMemoryStorage) StoreTask(task *Task) (int, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	if task.id != nil {
		if _, ok := storage.tasks[*task.id]; ok {
			// valid task id - overwrite
			storage.tasks[*task.id] = *task
			return *task.id, nil
		} else {
			return 0, errors.New("task with ID provided, but task does not exist")
		}
	}
	storage.lastId = storage.lastId + 1
	newId := storage.lastId
	storage.tasks[newId] = *task
	task.id = &newId
	return newId, nil
}

func (storage *InMemoryStorage) RetrieveTaskById(id int) (*Task, error) {
	storage.mu.RLock()
	if value, ok := storage.tasks[id]; ok {
		return &value, nil
	}
	defer storage.mu.RUnlock()
	return nil, errors.New("no task found with id " + string(id))
}

func CreateInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{tasks: make(map[int]Task), lastId: 0}
}
