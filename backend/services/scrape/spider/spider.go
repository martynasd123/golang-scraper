package spider

import (
	"log"
	"net/http"
	"net/url"
	"time"

	models "github.com/martynasd123/golang-scraper/models/scrape"
)

const MAX_INSTANCES = 5

type Spider struct {
	resultsChannel chan models.ProcessingUpdate
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
		return nil
	}
	return &models.LinkCrawledUpdate{
		Link:     link,
		Status:   resp.StatusCode,
		TimedOut: false,
	}
}

func (spider *Spider) startInstance() {
	for link := range spider.LinksChannel {
		result := spider.Crawl(link)
		if result != nil {
			select {
			case spider.resultsChannel <- result:
			default:
				// Channel is closed
			}
		}
	}
}

func (spider *Spider) Start() {
	for i := 0; i < MAX_INSTANCES; i++ {
		go spider.startInstance()
	}
}
