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
	UpdateChannel    chan ProcessingUpdate
	InterruptChannel chan struct{}
	link             *url.URL
}

type UpdatesSubscriber struct {
}

func CreateSeeker(link *url.URL) *Seeker {
	return &Seeker{
		UpdateChannel:    make(chan ProcessingUpdate),
		link:             link,
		InterruptChannel: make(chan struct{}, 1),
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
		select {
		case <-seeker.InterruptChannel:
			// Indicate that no more links will be sent
			close(spiderInstance.LinksChannel)
			// Wait for spider to finish
			<-done
			// Indicate interruption to the updates channel
			seeker.UpdateChannel <- &InterruptedUpdate{}
			return
		case spiderInstance.LinksChannel <- &link:
		}
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

	if err != nil || resp.StatusCode != http.StatusOK {
		seeker.UpdateChannel <- &ErrorUpdate{Error: fmt.Errorf("failed to GET initial page: %v", err)}
		return
	}
	defer closeHttp(resp)

	if seeker.checkInterrupt() {
		return
	}

	document, err := html.Parse(resp.Body)
	if err != nil {
		seeker.UpdateChannel <- &ErrorUpdate{Error: fmt.Errorf("failed to parse html: %v", err)}
		return
	}

	if seeker.checkInterrupt() {
		return
	}

	seeker.processPage(document)
}

func closeHttp(resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		log.Printf("failed to close response body: %v", err)
	}
}

// Checks if seeker is interrupted. Returns true if it is, and sends signal to update channel about the interruption
func (seeker *Seeker) checkInterrupt() bool {
	select {
	case <-seeker.InterruptChannel:
		seeker.UpdateChannel <- &InterruptedUpdate{}
		return true
	default:
		return false
	}
}
