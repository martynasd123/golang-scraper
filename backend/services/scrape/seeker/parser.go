package seeker

import (
	"log"
	"net/url"
	"strconv"
	"strings"

	models "github.com/martynasd123/golang-scraper/models/scrape"
	datatype "github.com/martynasd123/golang-scraper/utils/datatype"
	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) (*html.Attribute, bool) {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return &attr, true
		}
	}
	return nil, false
}

func extractLink(node *html.Node) (*url.URL, bool) {
	if node.Type == html.ElementNode && node.Data == "a" {
		attr, found := getAttr(node, "href")
		if found {
			url, err := url.Parse(attr.Val)
			if err == nil {
				if url.Scheme != "https" && url.Scheme != "http" && url.Scheme != "" {
					// Ignore mailto: and other similar links
					// Note that when scheme is empty, it means we encountered a relative link.
					// We want to allow those.
					return nil, false
				}
				url.Fragment = "" // Ignore fragment parts of links
			} else {
				log.Printf("Failed to parse link: %s. Got error: %s", attr, err)
			}
			return url, err == nil
		}
	}
	return nil, false
}

func traverse(rootNode *html.Node, handler func(node *html.Node) bool) {
	var nodesToProcess []*html.Node

	nodesToProcess = append(nodesToProcess, rootNode)

	var currentDepth = 0
	for len(nodesToProcess) != 0 {
		node := nodesToProcess[0]
		if handler(node) {
			break
		}
		nodesToProcess = nodesToProcess[1:]

		currentDepth = currentDepth + 1
		// Add all child nodes to be processed
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			nodesToProcess = append(nodesToProcess, c)
		}
	}
}

func isLoginForm(node *html.Node) bool {
	if node.Type != html.ElementNode || node.Data != "form" {
		return false
	}
	isLoginForm := false
	traverse(node, func(node *html.Node) bool {
		if node.Type == html.ElementNode && node.Data == "input" {
			attr, found := getAttr(node, "type")
			if found && attr.Val == "password" {
				isLoginForm = true
				return true // Make sure to break the loop
			}
		}
		return false
	})
	return isLoginForm
}

func isDoctypeNode(node *html.Node) bool {
	return node.Type == html.DoctypeNode
}

func isTitleNode(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "title"
}

func isHeadingTag(node *html.Node) (bool, int) {
	if node.Type != html.ElementNode || !strings.HasPrefix(node.Data, "h") || len(node.Data) != 2 {
		return false, 0
	}
	i, err := strconv.Atoi(node.Data[1:])
	if err != nil || i < 1 || i > 6 {
		return false, 0
	}
	return true, i
}

func parseHeaders(headerNode *html.Node) string {
	var title string
	traverse(headerNode, func(node *html.Node) bool {
		if isTitleNode(node) {
			title = node.FirstChild.Data
			return true
		}
		return false
	})
	return title
}

func isHeaderNode(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "head"
}

func (seeker *Seeker) isExternal(link string) bool {
	// todo
	return true
}

func calcInternalLinks(pageLink url.URL, links *datatype.Set[url.URL]) int {
	internalLinks := 0
	for _, link := range links.Values() {
		if link.Host == pageLink.Host {
			internalLinks = internalLinks + 1
		}
	}
	return internalLinks
}

func ParseBaseInfo(rootNode *html.Node, pageLink url.URL) *models.PageBaseInfo {
	var links = datatype.NewSet[url.URL]()
	var headingsByLevel [6]int
	var foundLoginForm = false
	var htmlVersion string
	var title string

	traverse(rootNode, func(node *html.Node) bool {
		if link, found := extractLink(node); found {
			links.Add(*pageLink.ResolveReference(link))
		}
		foundLoginForm = foundLoginForm || isLoginForm(node)
		if isDoctypeNode(node) {
			htmlVersion = strings.TrimPrefix(node.Data, "html")
		}
		if isHeaderNode(node) {
			title = parseHeaders(node)
		}
		if isHeading, level := isHeadingTag(node); isHeading {
			headingsByLevel[level] = headingsByLevel[level] + 1
		}
		return false
	})

	internalLinks := calcInternalLinks(pageLink, &links)

	return &models.PageBaseInfo{
		HtmlVersion:      htmlVersion,
		PageTitle:        title,
		LoginFormPresent: foundLoginForm,
		HeadingsByLevel:  headingsByLevel,
		InternalLinks:    internalLinks,
		ExternalLinks:    len(links) - internalLinks,
		Links:            links.Values(),
	}
}
