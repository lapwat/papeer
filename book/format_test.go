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

func TestToMarkdown(t *testing.T) {

	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	got := ToMarkdown(c)
	want := "Books\n=====\n\n- [Discours de la Méthode](https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/)clock 98 min read -\n1637\n\n- [The Twelve-Factor App](https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/)clock 22 min read -\n2011"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
	
}

func TestToEpub(t *testing.T) {

	filename := "Books.epub"
	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToEpub(c, "")

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
	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
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
	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})
	ToMobi(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}
	
}
