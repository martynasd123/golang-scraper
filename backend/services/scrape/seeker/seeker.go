package seeker

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	. "github.com/martynasd123/golang-scraper/models/scrape"
	"github.com/martynasd123/golang-scraper/services/scrape/spider"
	"golang.org/x/net/html"
)

type Seeker struct {
	UpdateChannel chan ProcessingUpdate
	link          *url.URL
}

type UpdatesSubscriber struct {
}

func CreateSeeker(link *url.URL) *Seeker {
	return &Seeker{
		UpdateChannel: make(chan ProcessingUpdate),
		link:          link,
	}
}

func (seeker *Seeker) processPage(rootNode *html.Node) {
	// Perform initial parsing
	baseInfo := ParseBaseInfo(rootNode, *seeker.link)

	// Send the base page info
	seeker.UpdateChannel <- &PageBaseInfoUpdate{
		BaseInfo: baseInfo,
	}

	// Instantiate spiderInstance
	spiderInstance := spider.CreateSpider(seeker.UpdateChannel)

	done := spiderInstance.Start()
	for _, link := range baseInfo.Links {
		spiderInstance.LinksChannel <- &link
	}
	// Indicate we have no more links to process
	close(spiderInstance.LinksChannel)

	// Wait for spiders to finish
	<-done

	seeker.UpdateChannel <- &FinishedUpdate{}
}

// Seek starts the seeking process. Updates are sent through seeker.UpdateChannel until it is closed.
func (seeker *Seeker) Seek() {
	defer close(seeker.UpdateChannel)

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(seeker.link.String())
	if err != nil {
		seeker.UpdateChannel <- &ErrorUpdate{Error: fmt.Errorf("failed to GET initial page: %v", err)}
		return
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		seeker.UpdateChannel <- &ErrorUpdate{Error: fmt.Errorf("failed to parse html: %v", err)}
		return
	}
	err = resp.Body.Close()
	if err != nil {
		log.Printf("failed to close response body: %v", err)
	}

	seeker.processPage(document)
}
