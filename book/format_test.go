package book

import (
	"errors"
	"os"
	"testing"
)

func TestFilename(t *testing.T) {

	got := Filename("This is a chapter / book")
	want := "This_is_a_chapter__book"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToMarkdownString(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})

	got := ToMarkdownString(c)
	want := "Example Domain\n==============\n\nThis domain is for use in illustrative examples in documents. You may use this\ndomain in literature without prior coordination or asking for permission.\n\n[More information...](https://www.iana.org/domains/example)\n\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToMarkdownPrintURL(t *testing.T) {

	config := NewScrapeConfigQuiet()
	config.PrintURL = true

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := ToMarkdownString(c)
	want := "Example Domain\n==============\n\n_Source: https://example.com/_\n\nThis domain is for use in illustrative examples in documents. You may use this\ndomain in literature without prior coordination or asking for permission.\n\n[More information...](https://www.iana.org/domains/example)\n\n\n"

	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

}

func TestToMarkdown(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToMarkdown(c, "")

	filename := "Example_Domain.md"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMarkdownFilename(t *testing.T) {

	filename := "ebook.md"
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToMarkdown(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToHtmlString(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})

	got := ToHtmlString(c)
	want := "<h1>Example Domain</h1>\n<div>\n    \n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToHtmlPrintURL(t *testing.T) {

	config := NewScrapeConfigQuiet()
	config.PrintURL = true

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{config}, 0, func(index int, name string) {})

	got := ToHtmlString(c)
	want := "<h1>Example Domain</h1>\n<p><i>Source: https://example.com/</i></p>\n<div>\n    \n    <p>This domain is for use in illustrative examples in documents. You may use this\n    domain in literature without prior coordination or asking for permission.</p>\n    <p><a href=\"https://www.iana.org/domains/example\">More information...</a></p>\n</div>"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToHtml(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToHtml(c, "")

	filename := "Example_Domain.html"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToHtmlFilename(t *testing.T) {

	filename := "ebook.html"
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToHtml(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToEpub(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToEpub(c, "")

	filename := "Example_Domain.epub"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToEpubFilename(t *testing.T) {

	filename := "ebook.epub"
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToEpub(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMobi(t *testing.T) {

	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToMobi(c, "")

	filename := "Example_Domain.mobi"
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}

func TestToMobiFilename(t *testing.T) {

	filename := "ebook.mobi"
	c := NewChapterFromURL("https://example.com/", "", []*ScrapeConfig{NewScrapeConfigQuiet()}, 0, func(index int, name string) {})
	ToMobi(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}
