# Golang scraper

This is a web scraper application that takes a website URL as input and provides general information about the contents of the page. It is built with a React frontend for the UI and a Golang backend for the API.

This project is not meant to be used in production environments.

## Features

- **HTML Version**: Detects the HTML version of the webpage.
- **Page Title**: Retrieves the title of the webpage.
- **HTML Heading Tags Count**: Counts the number of HTML heading tags (h1-h6) on the page.
- **Internal and External Links**: Determines the number of internal and external links on the page.
- **Inaccessible Links**: Identifies links that return 4xx or 5xx status codes.
- **Login Form Detection**: Indicates whether the page contains a login form.
- **Real-time progress**: Thanks to SSE, scraping progress can be viewed in real-time.
- **Task interruptions**: Tasks can be interrupted mid-scraping.

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [npm](https://www.npmjs.com/)
- [Golang](https://golang.org/)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/your-username/web-scraper.git
```

2. Run front-end

```bash
cd frontend && npm install && npm run dev
```

3. Run back-end (in separate shell)

```bash
cd backend && go run main.go
```

4. The system can be accessed at ``http://localhost:3000``


## Possible future improvements

- Some configuration framework (like [viper](https://github.com/spf13/viper)) integration, so that JWT key and other configuration variables could be stored separately and securely.
- Events sent through ``/api/scrape/task/:id/listen`` should be throttled in some way.
- A real database should be integrated. Currently, this integration relies on some in-memory storage implementations, which is not a very scalable solution. This should be relatively simple to do though, as I've abstracted away the storage logic.
- When a task is interrupted, the system waits for existing requests to finish before fully transitioning task to its final state. This could be improved by forcibly closing existing http connections and terminating task immediately.
- The back-end returns error messages as simple strings. While this is fine for a project of this size, a more streamlined approach could be used by utilizing a consistent error response object and defined error codes.
- User should be able to have several refresh tokens along with device identifiers.
- Front-End tests.
- Front-end accessibility improvements.

## Lessons learned

The real-time scrape updates feature relies on a one publisher - many subscribers system that I wrote (see ``stateBroadcaster.go``). While it works, I think that a simple polling for status technique would have been a better choice here due to the added complexity on the whole system.
