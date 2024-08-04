package book

import (
	"errors"
	"flag"
	"os"
	"strings"
	"testing"
)

func TestToManyMarkdowns(t *testing.T) {

	config := NewScrapeConfigsWikipedia()
	var directory string

	flag.StringVar(&directory, "directory", "MarkdownFiles", "MarkdownFiles")
	flag.Parse()

	c := NewChapterFromURL("https://atomicdesign.bradfrost.com/table-of-contents/", "", config, 0, func(index int, name string) {})
	flist := SaveChapterAndSubChaptersAsMarkdown(c, directory)

	got := len(flist)
	want := 9

	//check subchapters count
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}

	//check file path exist
	for _, file := range flist {
		if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
			t.Errorf("%s does not exist: %v", file, err)
		}
	}

}

func TestToCheckNameOfMarkdowns(t *testing.T) {

	config := NewScrapeConfigsWikipedia()
	var directory string
	containsSubChapter := false

	flag.StringVar(&directory, "directory", "MarkdownFiles", "MarkdownFiles")
	flag.Parse()

	c := NewChapterFromURL("https://atomicdesign.bradfrost.com/table-of-contents/", "", config, 0, func(index int, name string) {})
	sc := c.SubChapters()
	flist := SaveChapterAndSubChaptersAsMarkdown(c, directory)

	//check the file name contains the title of the web page
	for _, file := range flist {
		for _, subChapter := range sc {
			if strings.Contains(file, subChapter.name) {
				containsSubChapter = true
				break
			}
		}
		if !containsSubChapter {
			t.Errorf("File '%s' does not contain any subchapter names\n", file)
		}
	}

}
