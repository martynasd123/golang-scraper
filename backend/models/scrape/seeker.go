package scrape

import "net/url"

const (
	UPDATE_TYPE_PAGE_BASE_INFO = iota
	UPDATE_TYPE_LINK_CRAWLED
)

type PageBaseInfo struct {
	HtmlVersion      string
	PageTitle        string
	LoginFormPresent bool
	HeadingsByLevel  [6]int
	InternalLinks    int
	ExternalLinks    int
	Links            []url.URL
}

type ProcessingUpdate interface {
	Type() int
}

type PageBaseInfoUpdate struct {
	BaseInfo *PageBaseInfo
}

func (PageBaseInfoUpdate) Type() int {
	return UPDATE_TYPE_PAGE_BASE_INFO
}

type LinkCrawledUpdate struct {
	Link     *url.URL
	Status   int
	TimedOut bool
}

func (LinkCrawledUpdate) Type() int {
	return UPDATE_TYPE_LINK_CRAWLED
}

type ProcessingState struct {
	HtmlVersion       string
	PageTitle         string
	HtmlHeadingTags   [6]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	ContainsLoginForm bool
}
