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

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	got := ToMarkdownString(c)
	want := "Books\n=====\n\n- [Discours de la Méthode](https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/)clock 98 min read -\n   1637\n\n- [The Twelve-Factor App](https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/)clock 22 min read -\n   2011\n\n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToMarkdown(t *testing.T) {

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMarkdown(c, "")

	filename := "Books.md"
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
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
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

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	got := ToHtmlString(c)
	want := "<h1>Books</h1>\n    \n    <main>\n      \n  <article>\n    \n    <ul>\n      \n        <li>\n          <a href=\"https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/\">Discours de la Méthode</a>\n          <span>\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 98 min read -\n            1637\n          </span>\n        </li>\n      \n        <li>\n          <a href=\"https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/\">The Twelve-Factor App</a>\n          <span>\n            <svg xmlns=\"http://www.w3.org/2000/svg\" width=\"24\" height=\"24\" viewBox=\"0 0 24 24\" fill=\"none\" stroke=\"currentColor\" stroke-width=\"2\" stroke-linecap=\"round\" stroke-linejoin=\"round\">\n  <title>clock</title>\n  <circle cx=\"12\" cy=\"12\" r=\"10\"></circle><polyline points=\"12 6 12 12 16 14\"></polyline>\n</svg> 22 min read -\n            2011\n          </span>\n        </li>\n      \n    </ul>\n    \n\n\n\n  </article>\n\n    </main>\n  \n\n"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func TestToHtml(t *testing.T) {

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToHtml(c, "")

	filename := "Books.html"
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
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
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

	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToEpub(c, "")

	filename := "Books.epub"
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
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
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

	filename := "ebook.mobi"
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMobi(c, filename)

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
	c := NewChapterFromURL("https://books.lapw.at/", "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMobi(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}

}
