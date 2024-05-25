package seeker

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	models "github.com/martynasd123/golang-scraper/models/scrape"
	"github.com/martynasd123/golang-scraper/services/scrape/spider"
	"golang.org/x/net/html"
)

const MAX_PROCESSING_NODES = 100

type Seeker struct {
	link    *url.URL
	updates chan models.ProcessingUpdate
}

func CreateSeeker(link *url.URL) *Seeker {
	return &Seeker{link, make(chan models.ProcessingUpdate)}
}

func (seeker *Seeker) RegisterListener() (int, chan models.ProcessingState) {
	// todo
	return 0, make(chan models.ProcessingState)
}

func (seeker *Seeker) UnregisterListener(listenerId int) error {
	// todo
	return nil
}

func (seeker *Seeker) processPage(rootNode *html.Node) error {
	// Perform initial parsing
	baseInfo := ParseBaseInfo(rootNode, *seeker.link)

	// Send the base page info
	seeker.updates <- models.PageBaseInfoUpdate{
		BaseInfo: baseInfo,
	}

	// Instantiate spider
	spider := spider.CreateSpider(seeker.updates)

	spider.Start()
	for _, link := range baseInfo.Links {
		spider.LinksChannel <- &link
	}
	close(spider.LinksChannel)
	return nil
}

func (seeker *Seeker) monitor() {
	for update := range seeker.updates {
		fmt.Println(update)
	}
}

func (seeker *Seeker) Seek() {
	go seeker.monitor()

	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(seeker.link.String())
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	document, err := html.Parse(resp.Body)
	seeker.processPage(document)
}
