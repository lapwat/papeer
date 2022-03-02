package book

import (
	"testing"
	"time"
)

func TestBody(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Body()
	want := "<!doctype html>\n<html lang=\"en-us\">\n  <head>\n    <title>Books</title>\n    <link rel=\"shortcut icon\" href=\"/favicon.ico\" />\n    <meta charset=\"utf-8\" />\n    <meta name=\"generator\" content=\"Hugo 0.59.1\" />\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />\n    <meta name=\"author\" content=\"John Doe\" />\n    <meta name=\"description\" content=\" \" />\n    <link rel=\"stylesheet\" href=\"https://books.lapw.at/css/main.min.88e7083eff65effb7485b6e6f38d10afbec25093a6fac42d734ce9024d3defbd.css\" />\n\n    \n    <meta name=\"twitter:card\" content=\"summary\"/>\n<meta name=\"twitter:title\" content=\"Books\"/>\n<meta name=\"twitter:description\" content=\" \"/>\n\n    <meta property=\"og:title\" content=\"Books\" />\n<meta property=\"og:description\" content=\" \" />\n<meta property=\"og:type\" content=\"website\" />\n<meta property=\"og:url\" content=\"https://books.lapw.at/\" />\n\n\n\n  </head>\n  <body>\n    <header class=\"app-header\">\n      <a href=\"https://books.lapw.at/\"><img class=\"app-header-avatar\" src=\"/book.svg\" alt=\"John Doe\" /></a>\n      <h1>Books</h1>\n      <p> </p>\n      <div class=\"app-header-social\">\n        \n      </div>\n    </header>\n    <main class=\"app-container\">\n      \n  <article>\n    <h1>Books</h1>\n    <ul class=\"posts-list\">\n      \n        <li class=\"posts-list-item\">\n          <a class=\"posts-list-item-title\" href=\"https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/\">Discours de la Méthode</a>\n          <span class=\"posts-list-item-description\">\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-clock\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 98 min read -\n            1637\n          </span>\n        </li>\n      \n        <li class=\"posts-list-item\">\n          <a class=\"posts-list-item-title\" href=\"https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/\">The Twelve-Factor App</a>\n          <span class=\"posts-list-item-description\">\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"icon icon-clock\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 22 min read -\n            2011\n          </span>\n        </li>\n      \n    </ul>\n    \n\n\n\n  </article>\n\n    </main>\n  </body>\n</html>\n"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestName(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "Books"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestCustomName(t *testing.T) {

	config := NewScrapeConfig()
	config.UseLinkName = true
	c := NewChapterFromURL("https://books.lapw.at/", "Custom Name", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Name()
	want := "Custom Name"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestAuthor(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Author()
	want := "John Doe"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestContent(t *testing.T) {

	config := NewScrapeConfig()
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "<div id=\"readability-page-1\" class=\"page\">\n    \n    <main>\n      \n  <article>\n    \n    <ul>\n      \n        <li>\n          <a href=\"https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/\">Discours de la Méthode</a>\n          <span>\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 98 min read -\n            1637\n          </span>\n        </li>\n      \n        <li>\n          <a href=\"https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/\">The Twelve-Factor App</a>\n          <span>\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 22 min read -\n            2011\n          </span>\n        </li>\n      \n    </ul>\n    \n\n\n\n  </article>\n\n    </main>\n  \n\n</div>"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestDelay(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Delay = 500

	config1 := NewScrapeConfig()

	start := time.Now()
	NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})
	elapsed := time.Since(start)

	got := elapsed
	want := time.Duration(500) * time.Millisecond

	if got < want {
		t.Errorf("got %v, wanted min %v", got, want)
	}

}

func TestContentImagesOnly(t *testing.T) {

	config := NewScrapeConfig()
	config.ImagesOnly = true

	c := NewChapterFromURL("https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := "<img src=\"https://books.lapw.at/images/codebase-deploys.png\" alt=\"One codebase maps to many deploys\"/><img src=\"https://books.lapw.at/images/attached-resources.png\" alt=\"A production deploy attached to four backing services.\"/><img src=\"https://books.lapw.at/images/release.png\" alt=\"Code becomes a build, which is combined with config to create a release.\"/><img src=\"https://books.lapw.at/images/process-types.png\" alt=\"Scale is expressed as running processes, workload diversity is expressed as process types.\"/>"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChapters(t *testing.T) {

	config0 := NewScrapeConfig()
	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 2

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersRSS(t *testing.T) {

	config0 := NewScrapeConfig()
	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://blog.lapw.at/rss", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 8

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersSelector(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Selector = "section.concrete > article > h2 > a"

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://12factor.net/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 12

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimit(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Limit = 1

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 1

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestSubChaptersLimitOver(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Limit = 3

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := len(c.SubChapters())
	want := 2

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestReverse(t *testing.T) {

	config0 := NewScrapeConfig()
	config0.Reverse = true

	config1 := NewScrapeConfig()

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config0, config1}, 0, func(index int, name string) {})

	got := c.SubChapters()[0].Name()
	want := "The Twelve-Factor App"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestNotInclude(t *testing.T) {

	config := NewScrapeConfig()
	config.Include = false

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := c.Content()
	want := ""

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}
