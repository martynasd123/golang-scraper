package scrape

import "net/url"

const (
	UpdateTypePageBaseInfo = iota
	UpdateTypeLinkCrawled
	UpdateTypeError
	UpdateTypeFinished
	UpdateTypeInterrupted
)

const (
	// Did not start the task yet, but it has been queued
	StatusPending = "PENDING"
	// Performing the initial GET request and parsing the page
	StatusInitiating = "INITIATING"
	// Trying to access all related links
	StatusTryingLinks = "TRYING_LINKS"
	// Task was completed successfully
	StatusFinished = "FINISHED"
	// Task was interrupted
	StatusInterrupted = "INTERRUPTED"
	// Interrupt signal received - waiting for seeker/spiders to finish
	StatusInterrupting = "INTERRUPTING"
	// Error occurred
	StatusError = "ERROR"
)

// ProcessingUpdate is the interface for all seeker updates
type ProcessingUpdate interface {
	Type() int
}

// PageBaseInfo contains the initial information about the webpage.
type PageBaseInfo struct {
	HtmlVersion      string
	PageTitle        string
	LoginFormPresent bool
	HeadingsByLevel  [6]int
	InternalLinks    int
	ExternalLinks    int
	Links            []url.URL
}

type ErrorUpdate struct {
	Error error
}

func (e ErrorUpdate) Type() int {
	return UpdateTypeError
}

type InterruptedUpdate struct {
}

func (e InterruptedUpdate) Type() int {
	return UpdateTypeInterrupted
}

// PageBaseInfoUpdate is and update sent before the spiders are in action, after the initial GET request.
type PageBaseInfoUpdate struct {
	BaseInfo *PageBaseInfo
}

func (PageBaseInfoUpdate) Type() int {
	return UpdateTypePageBaseInfo
}

// FinishedUpdate is an update indicating that the scraping was completed successfully
type FinishedUpdate struct {
}

func (f FinishedUpdate) Type() int {
	return UpdateTypeFinished
}

// LinkCrawledUpdate is an update indicating that a link has been visited
type LinkCrawledUpdate struct {
	Link *url.URL
	// Http status of response
	Status int
	// Flag which indicates that a transport level error occurred
	TransportError bool
}

func (LinkCrawledUpdate) Type() int {
	return UpdateTypeLinkCrawled
}
