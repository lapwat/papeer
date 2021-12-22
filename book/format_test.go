package book

import (
	"errors"
	"os"
	"testing"
)

func TestToMarkdown(t *testing.T) {

	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()})

	got := ToMarkdown(c)
	want := "Books\n=====\n\n- [Discours de la Méthode](https://books.lapw.at/posts/ren%C3%A9-descartes-discours-de-la-m%C3%A9thode/)clock 98 min read -\n1637\n\n- [The Twelve-Factor App](https://books.lapw.at/posts/adam-wiggins-the-twelve-factor-app/)clock 22 min read -\n2011"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestToEpub(t *testing.T) {

	filename := "ebook.epub"
	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()})
	ToEpub(c, filename)

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}
}

func TestToEpubNoFilename(t *testing.T) {

	filename := "Books.epub"
	c := NewChapterFromURL("https://books.lapw.at/", []*ScrapeConfig{NewScrapeConfig()})
	ToEpub(c, "")

	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		t.Errorf("%s does not exist: %v", filename, err)
	} else {
		if err := os.Remove(filename); err != nil {
			t.Errorf("cannot remove %v: %v", filename, err)
		}
	}
}
