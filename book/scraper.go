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
)

func NewBookFromURL(url, selector string, recursive, include, images bool, limit, offset, delay, threads int) book {
	if recursive {
		chapters := tableOfContent(url, selector, limit, offset, delay, threads, include, images)

		b := New(chapters[0].Name(), chapters[0].Author())
		for _, c := range chapters {
			b.AddChapter(c)
		}

		return b
	} else {
		c := NewChapterFromURL(url, images)
		b := New(c.Name(), c.Author())
		b.AddChapter(c)
		return b
	}
}

func NewChapterFromURL(url string, images bool) chapter {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	content := strings.ReplaceAll(article.Content, "\n", "")

	if images {
		// parse html content
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err != nil {
			log.Fatal(err)
		}

		// extract images only
		content = ""
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			newContent, _ := goquery.OuterHtml(s)
			content += newContent
		})
	}

	return chapter{article.Title, article.Byline, content}
}

func tableOfContent(url, selector string, limit, offset, delay, threads int, include, images bool) []chapter {
	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	links, err := GetLinks(base, selector, limit, offset, include)
	if err != nil {
		log.Fatal(err)
	}

	chapters := make([]chapter, len(links))
	progress := NewProgress(links)

	if delay >= 0 {
		// synchronous mode

		for index, link := range links {
			// and then use it to parse relative URLs
			u, err := base.Parse(link.href)
			if err != nil {
				log.Fatal(err)
			}

			chapters[index] = NewChapterFromURL(u.String(), images)
			progress.Incr(index)

			// short sleep for last chapter to let the progress bar update
			if index == len(links)-1 {
				delay = 100
			}

			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

	} else {
		// asynchronous mode
		var wg sync.WaitGroup

		if threads == -1 {
			threads = len(links)
		}
		semaphore := make(chan bool, threads)

		for index, l := range links {

			wg.Add(1)
			semaphore <- true

			go func(index int, l link) {
				defer wg.Done()

				// and then use it to parse relative URLs
				u, err := base.Parse(l.href)
				if err != nil {
					log.Fatal(err)
				}

				chapters[index] = NewChapterFromURL(u.String(), images)
				progress.Incr(index)

				<-semaphore
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

func GetLinks(url *urllib.URL, selector string, limit, offset int, include bool) ([]link, error) {
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
		key := path

		// include element class in key if selector is set
		if !selectorSet {
			class := e.Attr("class")
			key = fmt.Sprintf("%s.%s", path, class)
		}

		if selectorSet || text != "" {
			pathLinks[key] = append(pathLinks[key], NewLink(href, text))
			pathCount[key] += len(text)

			if pathCount[key] > pathCount[pathMax] {
				pathMax = key
			}
		}
	})
	c.Visit(url.String())

	links := pathLinks[pathMax]
	if len(links) == 0 {
		return []link{}, fmt.Errorf("no link found for selector: %s", selector)
	}

	end := len(links)
	if limit != -1 {
		end = int(math.Min(float64(limit+offset), float64(len(links))))
	}

	links = links[offset:end]

	if include {
		c := NewChapterFromURL(url.String(), false)
		l := NewLink(url.String(), c.Name())
		links = append([]link{l}, links...)
	}

	return links, nil
}
