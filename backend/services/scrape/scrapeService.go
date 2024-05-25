package scrape

import (
	"errors"
	"net/url"

	models "github.com/martynasd123/golang-scraper/models/scrape"
	seeker "github.com/martynasd123/golang-scraper/services/scrape/seeker"
	storage "github.com/martynasd123/golang-scraper/services/scrape/storage"
)

type ScrapeService struct {
	// Task ID to seeker map
	seekerMap map[int]seeker.Seeker
	// Task storage interface
	storage   storage.TaskStorage
}

func NewScrapeService(storage storage.TaskStorage) *ScrapeService {
	return &ScrapeService{storage: storage, seekerMap: make(map[int]seeker.Seeker)}
}

func (service *ScrapeService) RegisterListener(seekerId int) (int, chan models.ProcessingState, error) {
	seeker, ok := service.seekerMap[seekerId]
	if !ok {
		return 0, nil, errors.New("invalid seeker id")
	}
	listenerId, channel := seeker.RegisterListener()
	return listenerId, channel, nil
}

func (service *ScrapeService) UnregisterListener(listenerId int, seekerId int) error {
	seeker, ok := service.seekerMap[seekerId]
	if !ok {
		return errors.New("invalid seeker id")
	}
	return seeker.UnregisterListener(listenerId)
}

// Creates a seeker for a given URL. Returns a unique ID associated with it
//
// Parameters:
//
//	link (url): The URL to scrape
//
// Returns:
//
//	int: The unique seeker identifier
func (ss *ScrapeService) AddTask(link *url.URL) (int, error) {
	task := storage.CreateTask(storage.TASK_STATUS_QUEUED, link)

	newId, err := ss.storage.StoreTask(task)
	if err != nil {
		return 0, err
	}

	seeker := seeker.CreateSeeker(link)

	ss.seekerMap[newId] = *seeker
	go seeker.Seek()
	return newId, nil
}
