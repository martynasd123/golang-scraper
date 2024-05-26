package scrape

import (
	"errors"
	"github.com/martynasd123/golang-scraper/models/scrape"
	. "github.com/martynasd123/golang-scraper/services/scrape/seeker"
	"github.com/martynasd123/golang-scraper/storage"
	"github.com/martynasd123/golang-scraper/utils/event"
	"log"
	"net/url"
)

const MaxInstances = 3

type ScrapeService struct {
	stateBroker *event.StateBroker[int, storage.Task]
	// Task ID to seeker map
	queuedTasks chan *storage.Task
	// Task storage interface
	storage storage.TaskDao
}

func CreateTaskService(taskStorage storage.TaskDao) *ScrapeService {
	scrapeService := &ScrapeService{
		storage:     taskStorage,
		stateBroker: event.CreateStateBroker[int, storage.Task](),
		queuedTasks: make(chan *storage.Task),
	}
	scrapeService.init()
	return scrapeService
}

func (service *ScrapeService) scrape() {
	for task := range service.queuedTasks {
		// Change task status and store in the database
		task.Status = scrape.StatusStarted
		// Persist in storage
		_, err := service.storage.StoreTask(task)
		if err != nil {
			log.Printf("could not store task: %v", err)
			continue
		}

		broadcaster, err := service.stateBroker.GetStateBroadcaster(*task.Id)
		if err != nil {
			log.Printf("could not retrieve state broker: %v", err)
			continue
		}

		// Notify subscribers of status started
		broadcaster.Publish(*task)

		seeker := CreateSeeker(&task.Link)
		go seeker.Seek()

		for update := range seeker.UpdateChannel {
			if update.Type() == scrape.UpdateTypePageBaseInfo {
				updateTaskBaseInfo(task, update.(*scrape.PageBaseInfoUpdate))
			} else if update.Type() == scrape.UpdateTypeLinkCrawled {
				handleLinkCrawled(task, update.(*scrape.LinkCrawledUpdate))
			} else if update.Type() == scrape.UpdateTypeError {
				handleError(task, update.(*scrape.ErrorUpdate))
				// Expecting this channel to close before next iteration
			} else if update.Type() == scrape.UpdateTypeFinished {
				handleFinished(task, update.(*scrape.FinishedUpdate))
				// Expecting this channel to close before next iteration
			} else {
				log.Fatalf("unsupported update type %d", update.Type())
			}
			broadcaster.Publish(*task)
		}

		// Update task status in the db
		_, err = service.storage.StoreTask(task)
		if err != nil {
			log.Printf("could not store task: %v", err)
		}

		// Delete and end state broker, as there will be no further updates to this task
		broadcaster.End()
		_ = service.stateBroker.DeleteStateBroadcaster(*task.Id)
	}
}

//goland:noinspection GoUnusedParameter
func handleFinished(task *storage.Task, update *scrape.FinishedUpdate) {
	task.Status = scrape.StatusFinished
}

func handleError(task *storage.Task, update *scrape.ErrorUpdate) {
	task.Status = scrape.StatusError
	err := update.Error.Error()
	task.Error = &err
}

func isInvalidHttpStatus(status int) bool {
	return status >= 400 && status < 600
}

func handleLinkCrawled(task *storage.Task, update *scrape.LinkCrawledUpdate) {
	if update.TransportError || isInvalidHttpStatus(update.Status) {
		*task.InaccessibleLinks = *task.InaccessibleLinks + 1
	}
	task.CrawledLinks = task.CrawledLinks + 1
}

func updateTaskBaseInfo(task *storage.Task, update *scrape.PageBaseInfoUpdate) {
	baseInfo := update.BaseInfo
	if baseInfo == nil {
		log.Fatalf("unexpected condition: BaseInfo nil in PageBaseInfoUpdate")
	}
	task.ExternalLinks = &baseInfo.ExternalLinks
	task.InternalLinks = &baseInfo.InternalLinks
	task.PageTitle = &baseInfo.PageTitle
	task.HtmlVersion = &baseInfo.HtmlVersion
	task.HeadingsByLevel = &baseInfo.HeadingsByLevel
	task.LoginFormPresent = &baseInfo.LoginFormPresent
	task.InaccessibleLinks = new(int)
}

func (service *ScrapeService) init() {
	for i := 0; i < MaxInstances; i++ {
		go service.scrape()
	}
}

func (service *ScrapeService) RegisterListener(taskId int) (err error, data <-chan storage.Task, done chan<- struct{}) {
	stateBroker, err := service.stateBroker.GetStateBroadcaster(taskId)
	if err != nil {
		task, err := service.storage.RetrieveTaskById(taskId)
		if err != nil {
			return err, nil, nil
		}
		if task.Status == scrape.StatusFinished || task.Status == scrape.StatusError {
			return errors.New("task already finished"), nil, nil
		}
		return errors.New("task not finished, but there is no state broker for it"), nil, nil
	}
	err = nil
	data, done = stateBroker.Listen()
	return
}

// AddTask Creates a seeker for a given URL. Returns a unique ID associated with it
//
// Parameters:
//
//	link (url): The URL to scrape
//
// Returns:
//
//	int: The unique seeker identifier
func (service *ScrapeService) AddTask(link *url.URL) (int, error) {
	task := storage.CreateTaskInitial(scrape.StatusPending, link)

	// Save the newly created task
	newId, err := service.storage.StoreTask(task)
	if err != nil {
		return 0, err
	}

	// CreateTaskService new state broadcaster
	broadcaster, err := service.stateBroker.AddStateBroadcaster(newId)
	if err != nil {
		return -1, err
	}

	// Publish initial state
	broadcaster.Start(*task)

	// Push to channel so that it starts processing
	go func() {
		service.queuedTasks <- task
	}()
	return newId, nil
}

func (service *ScrapeService) GetTaskById(id int) (*storage.Task, error) {
	return service.storage.RetrieveTaskById(id)
}

func (service *ScrapeService) GetAllTasks() []*storage.Task {
	return service.storage.GetAllTasks()
}
