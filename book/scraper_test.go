package book

import (
	"testing"
	"time"
)

func TestBody(t *testing.T) {

	config := NewScrapeConfigQuiet()
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Body()
	want := "<!doctype html><html lang=\"en\"><head><title>Example Domain</title><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><style>body{background:#eee;width:60vw;margin:15vh auto;font-family:system-ui,sans-serif}h1{font-size:1.5em}div{opacity:0.8}a:link,a:visited{color:#348}</style><body><div><h1>Example Domain</h1><p>This domain is for use in documentation examples without needing permission. Avoid use in operations.<p><a href=\"https://iana.org/domains/example\">Learn more</a></div></body></html>\n"

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
	want := "<div><p>This domain is for use in documentation examples without needing permission. Avoid use in operations.</p><p><a href=\"https://iana.org/domains/example\">Learn more</a></p></div>"

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
