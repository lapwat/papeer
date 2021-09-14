package book

import (
	"fmt"
	"log"
	urllib "net/url"
	"strings"
	"sync"
	"time"

	readability "github.com/go-shiori/go-readability"
	colly "github.com/gocolly/colly/v2"
	"github.com/gosuri/uiprogress"
)

type scraper struct {
	url string
}

func NewBookFromURL(url, selector string, include bool, delay int) Book {
	home := NewChapterFromURL(url)
	b := New(home.Name(), home.Author())

	chapters := tableOfContent(url, selector, delay)
	if include {
		b.AddChapter(home)
	}
	for _, c := range chapters {
		b.AddChapter(c)
	}

	return b
}

func NewChapterFromURL(url string) chapter {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	// metadata := fmt.Sprintf("URL     : %s\nTitle   : %s\nAuthor  : %s\nLength  : %d\nExcerpt : %s\nSiteName: %s\nImage   : %s\nFavicon : %s", url, article.Title, article.Byline, article.Length, article.Excerpt, article.SiteName, article.Image, article.Favicon)
	// fmt.Println(metadata)

	return chapter{article.Title, article.Byline, article.Content}
}

func tableOfContent(url, selector string, delay int) []chapter {
	c := colly.NewCollector()

	classesLinks := map[string][]map[string]string{}
	classesCount := map[string]int{}
	classMax := ""

	if selector == "" {
		c.OnHTML("a", func(e *colly.HTMLElement) {
			href := e.Attr("href")
			text := strings.TrimSpace(e.Text)
			class := e.Attr("class")

			//if class != "" && text != "" {
				classesLinks[class] = append(classesLinks[class], map[string]string{
					"href": href,
					"text": text,
				})

				classesCount[class]++

				if classesCount[class] > classesCount[classMax] {
					classMax = class
				}
			//}

		})
	} else {
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			href := e.Attr("href")
			text := strings.TrimSpace(e.Text)
			class := e.Attr("class")

			//if class != "" && text != "" {
				classesLinks[class] = append(classesLinks[class], map[string]string{
					"href": href,
					"text": text,
				})

				classesCount[class]++

				if classesCount[class] > classesCount[classMax] {
					classMax = class
				}
			//}
		})
	}
	c.Visit(url)
	fmt.Println(classesCount)
	links := classesLinks[classMax]

	chapters := make([]chapter, len(links))
	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	// init global progress bar
	uiprogress.Start()
	barGlobal := uiprogress.AddBar(len(links)).AppendCompleted().PrependElapsed()
	barGlobal.AppendFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Status: %d out of %d chapters", b.Current(), len(links))
	})

	// init progress bars
	bars := []*uiprogress.Bar{}
	for index, link := range links {
		bar := uiprogress.AddBar(1).AppendCompleted().PrependElapsed()
		barText := fmt.Sprintf("%d. %s", index+1, link["text"])
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return barText
		})
		bars = append(bars, bar)
	}

	if delay >= 0 {
		for index, link := range links {
			// and then use it to parse relative URLs
			u, err := base.Parse(link["href"])
			if err != nil {
				log.Fatal(err)
			}

			chapters[index] = NewChapterFromURL(u.String())

			bars[index].Incr()
			barGlobal.Incr()

			// do not wait after downloading last chapter
			if index < len(links)-1 {
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}

		}

	} else {
		var wg sync.WaitGroup
		for index, link := range links {

			wg.Add(1)
			go func(index int, link map[string]string) {
				defer wg.Done()

				// and then use it to parse relative URLs
				u, err := base.Parse(link["href"])
				if err != nil {
					log.Fatal(err)
				}

				chapters[index] = NewChapterFromURL(u.String())

				bars[index].Incr()
				barGlobal.Incr()
			}(index, link)
		}
		wg.Wait()
	}
	return chapters
}
