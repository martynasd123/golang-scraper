package seeker

import (
	"github.com/martynasd123/golang-scraper/models/scrape"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"testing"
)

func TestParseBaseInfo(t *testing.T) {
	testHTML := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Test Page</title>
	</head>
	<body>
		<h1>Welcome</h1>
		<a href="https://example.com">Link</a>
		<form>
			<input type="text" name="username">
			<input type="password" name="password">
			<button type="submit">Submit</button>
		</form>
	</body>
	</html>
	`
	testReader := strings.NewReader(testHTML)
	testNode, err := html.Parse(testReader)
	require.NoError(t, err)

	expectedBaseInfo := &scrape.PageBaseInfo{
		HtmlVersion:      "HTML 5.0",
		PageTitle:        "Test Page",
		LoginFormPresent: true,
		HeadingsByLevel:  [6]int{0, 1, 0, 0, 0, 0},
		InternalLinks:    1,
		ExternalLinks:    0,
		Links: []url.URL{
			{Scheme: "https", Host: "example.com", Path: ""},
		},
	}

	pageURL, _ := url.Parse("https://example.com")
	resultBaseInfo := ParseBaseInfo(testNode, *pageURL)

	require.Equal(t, expectedBaseInfo, resultBaseInfo)
}

func TestParseBaseInfo_NoLoginFormAndExternalLink(t *testing.T) {
	testHTML := `
	<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
	<html>
	<head>
		<title>Test Page</title>
	</head>
	<body>
		<h4>Welcome</h4>
		<a href="https://example.com">Link</a>
		<a href="https://other-site.com">Link</a>
	</body>
	</html>
	`
	testReader := strings.NewReader(testHTML)
	testNode, err := html.Parse(testReader)
	require.NoError(t, err)

	expectedBaseInfo := &scrape.PageBaseInfo{
		HtmlVersion:      "HTML 4.01",
		PageTitle:        "Test Page",
		LoginFormPresent: false,
		HeadingsByLevel:  [6]int{0, 0, 0, 0, 1, 0},
		InternalLinks:    1,
		ExternalLinks:    1,
		Links: []url.URL{
			{Scheme: "https", Host: "example.com", Path: ""},
			{Scheme: "https", Host: "other-site.com", Path: ""},
		},
	}

	pageURL, err := url.Parse("https://example.com")
	require.NoError(t, err)

	resultBaseInfo := ParseBaseInfo(testNode, *pageURL)

	require.Equal(t, expectedBaseInfo, resultBaseInfo)
}
