package book

import (
	"testing"
	"time"
)

func TestBody(t *testing.T) {

	config := NewScrapeConfigQuiet()
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Body()
	want := "<!doctype html>\n<html>\n<head>\n    <title>Example Domain</title>\n\n    <meta charset=\"utf-8\" />\n    <meta http-equiv=\"Content-type\" content=\"text/html; charset=utf-8\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <style type=\"text/css\">\n    body {\n        background-color: #f0f0f2;\n        margin: 0;\n        padding: 0;\n        font-family: -apple-system, system-ui, BlinkMacSystemFont, \"Segoe UI\", \"Open Sans\", \"Helvetica Neue\", Helvetica, Arial, sans-serif;\n        \n    }\n    div {\n        width: 600px;\n        margin: 5em auto;\n        padding: 2em;\n        background-color: #fdfdff;\n        border-radius: 0.5em;\n        box-shadow: 2px 3px 7px 2px rgba(0,0,0,0.02);\n    }\n    a:link, a:visited {\n        color: #38488f;\n        text-decoration: none;\n    }\n    @media (max-width: 700px) {\n        div {\n            margin: 0 auto;\n            width: auto;\n        }\n    }\n    </style>    \n</head>\n\n<body>\n<div>\n    <h1>Example Domain</h1>\n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>\n</body>\n</html>\n"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestName(t *testing.T) {

	config := NewScrapeConfigQuiet()
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "Example Domain"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestCustomName(t *testing.T) {

	config := NewScrapeConfigQuiet()
	config.UseLinkName = true
	c := NewChapterFromURL("https://example.com/", "Custom Name", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "Custom Name"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestAuthor(t *testing.T) {

	config := NewScrapeConfigQuiet()
	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Author()
	want := "Adam Wiggins"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestContent(t *testing.T) {

	config := NewScrapeConfigQuiet()
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "<div>\n    \n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestDelay(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Delay = 500

	config1 := NewScrapeConfigQuiet()

	start := time.Now()
	NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})
	elapsed := time.Since(start)

	got := elapsed
	want := time.Duration(500) * time.Millisecond

	if got < want {
		t.Errorf("got %v, wanted min %v", got, want)
	}

}

func TestContentImagesOnly(t *testing.T) {

	config := NewScrapeConfigQuiet()
	config.ImagesOnly = true

	c := NewChapterFromURL("https://12factor.net/codebase", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "<img src=\"https://12factor.net/images/codebase-deploys.png\" alt=\"One codebase maps to many deploys\"/>"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChapters(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Selector = ".concrete>article>h2>a"
	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersRSS(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://blog.nginx.org/feed", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 10

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersSelector(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Selector = "section.concrete>article>h2>a"

	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimit(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Selector = "section.concrete>article>h2>a"
	config0.Limit = 1

	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 1

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimitOver(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Selector = "section.concrete>article>h2>a"
	config0.Limit = 15

	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestReverse(t *testing.T) {

	config0 := NewScrapeConfigQuiet()
	config0.Selector = "section.concrete>article>h2>a"
	config0.Reverse = true

	config1 := NewScrapeConfigQuiet()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := c.SubChapters()[0].URL()
	want := "https://12factor.net/admin-processes"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestNotInclude(t *testing.T) {

	config := NewScrapeConfigQuiet()
	config.Include = false

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := ""

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}
