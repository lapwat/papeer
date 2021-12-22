package book

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	urllib "net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/go-shiori/go-readability"
	colly "github.com/gocolly/colly/v2"
)

type ScrapeConfig struct {
	selector   string
	limit      int
	include    bool
	imagesOnly bool
}

func NewScrapeConfig() *ScrapeConfig {
	return &ScrapeConfig{"", -1, true, false}
}

func NewBookFromURL(url, selector, name, author string, recursive, include, imagesOnly bool, limit, offset, delay, threads int) book {
	config1 := NewScrapeConfig()
	config1.imagesOnly = imagesOnly

	var chapters []chapter
	var home chapter

	if recursive {
		config2 := NewScrapeConfig()
		config2.selector = selector
		config2.limit = limit
		config2.include = include
		config2.imagesOnly = imagesOnly
		chapters, home = tableOfContent(url, config1.selector, config1.limit, offset, delay, threads, config1.include)
	} else {
		chapters = []chapter{NewChapterFromURL(url, []*ScrapeConfig{config1})}
		home = chapters[0]
	}

	if len(name) == 0 {
		name = home.Name()
	}

	if len(author) == 0 {
		author = home.Author()
	}

	b := New(name, author)
	for _, c := range chapters {
		b.AddChapter(c)
	}

	return b
}

func NewChapterFromURL(url string, configs []*ScrapeConfig) chapter {
	config := configs[0]
	content := ""

	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	subchapters := []chapter{}
	if len(configs) > 1 {
		// add subchapters

		links, _, err := GetLinks(base, config.selector, config.limit, 0, false)
		if err != nil {
			log.Fatal(err)
		}

		for _, link := range links {
			// and then use it to parse relative URLs
			u, err := base.Parse(link.href)
			if err != nil {
				log.Fatal(err)
			}

			subchapters = append(subchapters, NewChapterFromURL(u.String(), configs[1:]))
		}
	}

	// we want the metadata anyway

	// get page body
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// duplicate response stream
	readabilityReader := &bytes.Buffer{}
	bodyReader := io.TeeReader(response.Body, readabilityReader)

	// extract HTML body
	body, err := io.ReadAll(bodyReader)

	// extract content
	article, err := readability.FromReader(readabilityReader, base)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	// we don't care about the content if we do not include this level

	if config.include {
		content = article.Content

		// extract images
		if config.imagesOnly {

			// parse HTML
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
			if err != nil {
				log.Fatal(err)
			}

			// append every image to content
			content = ""
			doc.Find("img").Each(func(i int, s *goquery.Selection) {
				imageTag, _ := goquery.OuterHtml(s)
				content += imageTag
			})

		}
	}

	return chapter{string(body), article.Title, article.Byline, content, subchapters, config}
}

func tableOfContent(url, selector string, limit, offset, delay, threads int, include bool) ([]chapter, chapter) {
	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	links, home, err := GetLinks(base, selector, limit, offset, include)
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

			chapters[index] = NewChapterFromURL(u.String(), []*ScrapeConfig{NewScrapeConfig()})
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

				chapters[index] = NewChapterFromURL(u.String(), []*ScrapeConfig{NewScrapeConfig()})
				progress.Incr(index)

				<-semaphore
			}(index, l)
		}
		wg.Wait()
	}

	return chapters, home
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

func GetLinks(url *urllib.URL, selector string, limit, offset int, include bool) ([]link, chapter, error) {
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
		return []link{}, chapter{}, fmt.Errorf("no link found for selector: %s", selector)
	}

	end := len(links)
	if limit != -1 {
		end = int(math.Min(float64(limit+offset), float64(len(links))))
	}

	links = links[offset:end]

	home := NewChapterFromURL(url.String(), []*ScrapeConfig{NewScrapeConfig()})

	if include {
		l := NewLink(url.String(), home.Name())
		links = append([]link{l}, links...)
	}

	return links, home, nil
}
