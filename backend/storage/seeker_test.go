package storage

import (
	"github.com/martynasd123/golang-scraper/models/scrape"
	"net/url"
	"testing"
	"time"
)

func TestStoreTask(t *testing.T) {
	dao := CreateTaskInMemoryDao()

	link, _ := url.Parse("http://example.com")
	task := CreateTaskInitial(scrape.StatusPending, link, getSampleTime())

	id, err := dao.StoreTask(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected valid task id, got leq 0")
	}

	// Store the same task again, should overwrite
	task.Status = "completed"
	task.ExternalLinks = new(int)
	*task.ExternalLinks = 5

	id, err = dao.StoreTask(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if *task.Id != id {
		t.Fatalf("expected task id %v, got %v", *task.Id, id)
	}

	nonExistentTask := &Task{Id: new(int)}
	*nonExistentTask.Id = 999
	_, err = dao.StoreTask(nonExistentTask)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestRetrieveTaskById(t *testing.T) {
	dao := CreateTaskInMemoryDao()

	link, _ := url.Parse("http://example.com")
	task := CreateTaskInitial(scrape.StatusPending, link, getSampleTime())

	id, err := dao.StoreTask(task)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	retrievedTask, err := dao.RetrieveTaskById(id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if retrievedTask.Status != task.Status {
		t.Fatalf("expected status %v, got %v", task.Status, retrievedTask.Status)
	}

	_, err = dao.RetrieveTaskById(999)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestGetAllTasks(t *testing.T) {
	dao := CreateTaskInMemoryDao()

	link1, _ := url.Parse("http://example1.com")
	task1 := CreateTaskInitial("pending", link1, getSampleTime())

	link2, _ := url.Parse("http://example2.com")
	task2 := CreateTaskInitial("completed", link2, getSampleTime().Add(time.Second))

	_, err := dao.StoreTask(task1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Delay to ensure tasks have different creation times
	time.Sleep(1 * time.Second)

	_, err = dao.StoreTask(task2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	tasks := dao.GetAllTasks()
	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %v", len(tasks))
	}

	if tasks[1].CTime.After(tasks[0].CTime) {
		t.Fatalf("expected tasks to be sorted by creation time descending")
	}
}

func getSampleTime() time.Time {
	parsedTime, err := time.Parse("2006-01-02 15:04:05", "2023-05-27 14:23:45")
	if err != nil {
		panic(err)
	}
	return parsedTime
}
