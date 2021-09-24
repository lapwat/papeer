package book

import (
	"fmt"
	"log"
	"math"
	urllib "net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/go-shiori/go-readability"
	colly "github.com/gocolly/colly/v2"
	"github.com/gosuri/uiprogress"
)

func NewBookFromURL(url, selector string, recursive, include bool, limit, delay int) book {
	if recursive {
		home := NewChapterFromURL(url)
		b := New(home.Name(), home.Author())

		chapters := tableOfContent(url, selector, limit, delay)
		if include {
			b.AddChapter(home)
		}
		for _, c := range chapters {
			b.AddChapter(c)
		}

		return b
	} else {
		c := NewChapterFromURL(url)
		b := New(c.Name(), c.Author())
		b.AddChapter(c)
		return b
	}
}

func NewChapterFromURL(url string) chapter {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	return chapter{article.Title, article.Byline, article.Content}
}

func tableOfContent(url, selector string, limit, delay int) []chapter {
	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	links := GetLinks(base, selector)
	if limit != -1 {
		limit = int(math.Min(float64(limit), float64(len(links))))
		links = links[:limit]
	}

	chapters := make([]chapter, len(links))

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
		barText := fmt.Sprintf("%d. %s", index+1, link.text)
		bar.AppendFunc(func(b *uiprogress.Bar) string {
			return barText
		})
		bars = append(bars, bar)
	}

	if delay >= 0 {
		for index, link := range links {
			// and then use it to parse relative URLs
			u, err := base.Parse(link.href)
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
		for index, l := range links {

			wg.Add(1)
			go func(index int, l link) {
				defer wg.Done()

				// and then use it to parse relative URLs
				u, err := base.Parse(l.href)
				if err != nil {
					log.Fatal(err)
				}

				chapters[index] = NewChapterFromURL(u.String())

				bars[index].Incr()
				barGlobal.Incr()
			}(index, l)
		}
		wg.Wait()
	}
	return chapters
}

func GetPath(elm *goquery.Selection) string {
	path := []string{}

	for {
		selector := strings.ToLower(goquery.NodeName(elm))
		if selector == "" {
			break
		}

		path = append(path, selector)
		elm = elm.Parent()
	}

	join := strings.Join(path, "<")
	return join
}


func GetLinks(url *urllib.URL, selector string) []link {
	selectorSet := true
	if selector == "" {
		selector = "a"
		selectorSet = false
	}

	// visit and count link classes
	pathLinks := map[string][]link{}
	pathCount := map[string]int{}
	pathMax := ""

	c := colly.NewCollector()
	c.OnHTML(selector, func(e *colly.HTMLElement) {
		href := e.Attr("href")
		text := strings.TrimSpace(e.Text)
		path := GetPath(e.DOM)
		class := e.Attr("class")
		key := fmt.Sprintf("%s.%s", path, class)

		if selectorSet || text != "" {
			pathLinks[key] = append(pathLinks[key], NewLink(href, text, class))
			pathCount[key] += len(text)
			// pathCount[key]++

			if pathCount[key] > pathCount[pathMax] {
				pathMax = key
			}
		}
	})
	c.Visit(url.String())
	return pathLinks[pathMax]

	// // visit and count link classes
	// classesLinks := map[string][]link{}
	// classesCount := map[string]int{}
	// classMax := ""

	// c := colly.NewCollector()
	// c.OnHTML(selector, func(e *colly.HTMLElement) {
	// 	href := e.Attr("href")
	// 	text := strings.TrimSpace(e.Text)
	// 	class := e.Attr("class")

	// 	if selectorSet || class != "" && text != "" {
	// 		classesLinks[class] = append(classesLinks[class], NewLink(href, text))
	// 		classesCount[class]++

	// 		if classesCount[class] > classesCount[classMax] {
	// 			classMax = class
	// 		}
	// 	}
	// })
	// c.Visit(url.String())
	// return classesLinks[classMax]
}
