package scrape

import (
	"errors"
	"github.com/martynasd123/golang-scraper/models/scrape"
	. "github.com/martynasd123/golang-scraper/services/scrape/seeker"
	"github.com/martynasd123/golang-scraper/storage"
	"github.com/martynasd123/golang-scraper/utils/event"
	"log"
	"net/url"
	"sync"
)

var (
	ErrTaskInFinalState     = errors.New("task is already in final state")
	ErrInterruptAlreadySent = errors.New("interrupt signal already sent")
)

const MaxInstances = 3

type ScrapeService struct {
	// A map of task id to channel, which is to be used when interrupting the task
	interruptSignalMap map[int]chan<- struct{}
	// Locking this mutex locks task status transitions PENDING -> INITIATING
	interruptMu sync.Mutex
	stateBroker *event.StateBroker[int, storage.Task]
	// Task ID to seeker map
	queuedTasks chan int
	// Task storage interface
	storage storage.TaskDao
}

func CreateTaskService(taskStorage storage.TaskDao) *ScrapeService {
	scrapeService := &ScrapeService{
		storage:            taskStorage,
		stateBroker:        event.CreateStateBroker[int, storage.Task](),
		queuedTasks:        make(chan int),
		interruptSignalMap: make(map[int]chan<- struct{}),
		interruptMu:        sync.Mutex{},
	}
	scrapeService.init()
	return scrapeService
}

func (service *ScrapeService) setInterruptChannelForTask(taskId int, interruptChannel chan struct{}) {
	service.interruptMu.Lock()
	service.interruptSignalMap[taskId] = interruptChannel
	service.interruptMu.Unlock()
}

func (service *ScrapeService) deleteInterruptChannelForTask(taskId int) {
	service.interruptMu.Lock()
	delete(service.interruptSignalMap, taskId)
	service.interruptMu.Unlock()
}

func (service *ScrapeService) processTask(taskId int) {
	interruptChannel := make(chan struct{}, 1)

	// Set the channel, through which this task can be interrupted
	service.setInterruptChannelForTask(taskId, interruptChannel)
	defer service.deleteInterruptChannelForTask(taskId)

	broadcaster, err := service.stateBroker.GetStateBroadcaster(taskId)
	if err != nil {
		log.Printf("could not retrieve state broker: %v", err)
		return
	}
	defer service.destroyStateBroadcaster(taskId, broadcaster)

	task, err := service.storage.RetrieveTaskById(taskId)
	if err != nil {
		log.Printf("Error retrieving task %d: %v", taskId, err)
		return
	}
	if task.Status == scrape.StatusInterrupted {
		// Task has been interrupted before it started processing
		return
	}

	// Change task status and store in the database
	task.Status = scrape.StatusInitiating

	// Persist in storage
	_, err = service.storage.StoreTask(task)
	if err != nil {
		log.Printf("could not store task: %v", err)
		return
	}

	// Notify subscribers of status started
	broadcaster.Publish(*task)
	seeker := CreateSeeker(&task.Link)
	go seeker.Seek()

	func() {
		for {
			select {
			case <-interruptChannel:
				handleInterruptBegin(task)
				broadcaster.Publish(*task)
				seeker.InterruptChannel <- struct{}{}
			case update, ok := <-seeker.UpdateChannel:
				if !ok {
					return
				}
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
				} else if update.Type() == scrape.UpdateTypeInterrupted {
					// Update was interrupted
					handleInterruptFinish(task, update.(*scrape.InterruptedUpdate))
					// Expecting this channel to close before next iteration
				} else {
					log.Fatalf("unsupported update type %d", update.Type())
				}
				broadcaster.Publish(*task)
			}
		}
	}()

	// Update task status in the db
	_, err = service.storage.StoreTask(task)
	if err != nil {
		log.Printf("could not store task: %v", err)
	}
}

func handleInterruptBegin(task *storage.Task) {
	task.Status = scrape.StatusInterrupting
}

func (service *ScrapeService) destroyStateBroadcaster(taskId int, broadcaster *event.StateBroadcaster[storage.Task]) {
	broadcaster.End()
	err := service.stateBroker.DeleteStateBroadcaster(taskId)
	if err != nil {
		log.Printf("could not delete state broadcaster: %v", err)
		return
	}
}

func (service *ScrapeService) scrape() {
	for taskId := range service.queuedTasks {
		service.processTask(taskId)
	}
}

//goland:noinspection GoUnusedParameter
func handleFinished(task *storage.Task, update *scrape.FinishedUpdate) {
	task.Status = scrape.StatusFinished
}

//goland:noinspection GoUnusedParameter
func handleInterruptFinish(task *storage.Task, update *scrape.InterruptedUpdate) {
	task.Status = scrape.StatusInterrupted
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
	if task.Status == scrape.StatusInitiating {
		task.Status = scrape.StatusTryingLinks
	}
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
		if task.Status == scrape.StatusFinished || task.Status == scrape.StatusError || task.Status == scrape.StatusInterrupted {
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
		service.queuedTasks <- *task.Id
	}()
	return newId, nil
}

func (service *ScrapeService) InterruptTask(id int) error {
	service.interruptMu.Lock()
	defer service.interruptMu.Unlock()

	task, err := service.GetTaskById(id)
	if err != nil {
		return err
	}
	if task.Status == scrape.StatusPending {
		// In pending state - we can just update status, and it will not be picked up
		return service.interruptPendingTask(task)
	}

	interruptChannel, found := service.interruptSignalMap[id]
	if found {
		// Task is currently processing - need to send interrupt signal
		interruptChannel <- struct{}{}
		delete(service.interruptSignalMap, id)
		return nil
	}

	// Task is in a final state - return error
	switch task.Status {
	case scrape.StatusInterrupted, scrape.StatusError, scrape.StatusFinished:
		return ErrTaskInFinalState
	default:
		return ErrInterruptAlreadySent
	}
}

func (service *ScrapeService) interruptPendingTask(task *storage.Task) error {
	task.Status = scrape.StatusInterrupted
	_, err := service.storage.StoreTask(task)
	if err != nil {
		return err
	}
	broadcaster, err := service.stateBroker.GetStateBroadcaster(*task.Id)
	if err != nil {
		log.Printf("could not retrieve state broadcaster: %v", err)
		return err
	}
	// Publish update so that the subscribers know this task has been interrupted
	broadcaster.Publish(*task)
	return nil
}

func (service *ScrapeService) GetTaskById(id int) (*storage.Task, error) {
	return service.storage.RetrieveTaskById(id)
}

func (service *ScrapeService) GetAllTasks() []*storage.Task {
	return service.storage.GetAllTasks()
}
