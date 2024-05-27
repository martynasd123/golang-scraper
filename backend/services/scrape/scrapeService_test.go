package scrape_test

import (
	"fmt"
	scrapeStorage "github.com/martynasd123/golang-scraper/models/scrape"
	"github.com/martynasd123/golang-scraper/services/scrape"
	"github.com/martynasd123/golang-scraper/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestScrapeService_AddTaskAndListenForUpdates(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/error-response", errorResponseHandler)
	mux.HandleFunc("/success", createHtmlResponseHandler())
	mux.HandleFunc("/", createHtmlResponseHandler("/refuse-connection", "/error-response", "/success"))

	server := httptest.NewServer(mux)
	serverUrl, _ := url.Parse(server.URL)

	taskStorage := storage.CreateTaskInMemoryDao()

	service := scrape.CreateTaskService(taskStorage)

	_, data, _, err := service.AddTaskAndListenForUpdates(serverUrl)
	require.NoError(t, err)

	update := <-data
	assert.Equal(t, scrapeStorage.StatusPending, update.Status)

	update = <-data
	assert.Equal(t, scrapeStorage.StatusInitiating, update.Status)

	for range 4 {
		update = <-data
		assert.Equal(t, scrapeStorage.StatusTryingLinks, update.Status)
	}

	update = <-data
	assert.Equal(t, scrapeStorage.StatusFinished, update.Status)
	assert.Equal(t, 1, *update.InaccessibleLinks)
	assert.Equal(t, 3, update.CrawledLinks)
	assert.Equal(t, 3, *update.InternalLinks)
}

func TestScrapeService_AddTaskAndListenForUpdatesWhenPrimaryPageError(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", errorResponseHandler)

	server := httptest.NewServer(mux)
	serverUrl, _ := url.Parse(server.URL)

	taskStorage := storage.CreateTaskInMemoryDao()

	service := scrape.CreateTaskService(taskStorage)

	_, data, _, err := service.AddTaskAndListenForUpdates(serverUrl)
	require.NoError(t, err)

	update := <-data
	require.Equal(t, scrapeStorage.StatusPending, update.Status)

	update = <-data
	require.Equal(t, scrapeStorage.StatusInitiating, update.Status)

	update = <-data
	require.Equal(t, scrapeStorage.StatusError, update.Status)
}

func errorResponseHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func createHtmlResponseHandler(links ...string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html>")
		fmt.Fprintf(w, "<body>")
		for _, link := range links {
			fmt.Fprintf(w, "<a href=\"%s\">%s</a>", link, link)
		}
		fmt.Fprintf(w, "</body>")
		fmt.Fprintf(w, "</html>")
	}
}
