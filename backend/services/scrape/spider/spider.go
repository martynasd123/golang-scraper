package spider

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	models "github.com/martynasd123/golang-scraper/models/scrape"
)

const MaxInstances = 15

type Spider struct {
	resultsChannel chan models.ProcessingUpdate
	done           chan struct{}
	waitGroup      sync.WaitGroup
	LinksChannel   chan *url.URL
}

func CreateSpider(resultsChannel chan models.ProcessingUpdate) *Spider {
	return &Spider{
		resultsChannel: resultsChannel,
		LinksChannel:   make(chan *url.URL),
	}
}

func (spider *Spider) Crawl(link *url.URL) *models.LinkCrawledUpdate {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(link.String())
	if err != nil {
		log.Printf("error while crawling webpage link: %s. Got error: %s", link.String(), err)
		return &models.LinkCrawledUpdate{
			Link:           link,
			Status:         -1,
			TransportError: true,
		}
	}
	return &models.LinkCrawledUpdate{
		Link:           link,
		Status:         resp.StatusCode,
		TransportError: false,
	}
}

func (spider *Spider) startInstance() {
	for link := range spider.LinksChannel {
		result := spider.Crawl(link)
		if result != nil {
			spider.resultsChannel <- result
		}
	}
	spider.waitGroup.Done()
}

func (spider *Spider) Start() chan struct{} {
	spider.done = make(chan struct{})
	spider.waitGroup = sync.WaitGroup{}
	for i := 0; i < MaxInstances; i++ {
		spider.waitGroup.Add(1)
		go spider.startInstance()
	}
	go func() {
		spider.waitGroup.Wait()
		spider.done <- struct{}{}
		close(spider.done)
	}()
	return spider.done
}
