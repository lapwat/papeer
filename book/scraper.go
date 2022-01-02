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
	depth      int
	selector   string
	limit      int
	offset     int
	delay      int
	threads    int
	include    bool
	imagesOnly bool
}

func NewScrapeConfig() *ScrapeConfig {
	return &ScrapeConfig{0, "", -1, 0, -1, -1, true, false}
}

func NewScrapeConfigsAjin() []*ScrapeConfig {
	config0 := NewScrapeConfig()
	config0.depth = 0
	config0.selector = ".dt>a"
	config0.limit = 3
	config0.offset = 0
	config0.delay = 5000
	config0.include = false

	config1 := NewScrapeConfig()
	config1.depth = 1
	config1.selector = ".nav_apb>a"
	config1.limit = 3
	config1.offset = 1
	config1.delay = 5000
	config1.include = false

	config2 := NewScrapeConfig()
	config2.depth = 2
	config2.imagesOnly = true

	return []*ScrapeConfig{config0, config1, config2}
}

func NewScrapeConfigsWikipedia() []*ScrapeConfig {
	config0 := NewScrapeConfig()
	config0.depth = 0
	config0.threads = -1
	config0.include = true

	config1 := NewScrapeConfig()
	config1.depth = 1
	config1.include = true

	return []*ScrapeConfig{config0, config1}
}

func NewScrapeConfigFake() *ScrapeConfig {
	config := NewScrapeConfig()
	config.include = false

	return config
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
		config2.offset = offset
		config2.delay = delay
		config2.threads = threads
		config2.include = include
		config2.imagesOnly = imagesOnly
		chapters, home = tableOfContent(url, config2, config1)
	} else {
		chapters = []chapter{NewChapterFromURL(url, []*ScrapeConfig{config1}, 0, func(index int, name string) {})}
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

func NewChapterFromURL(url string, configs []*ScrapeConfig, index int, updateProgressBarName func(index int, name string)) chapter {
	config := configs[0]

	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

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

	// extract article content and metadata
	article, err := readability.FromReader(readabilityReader, base)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}
	name := article.Title

	// notify progress bar with new name
	updateProgressBarName(index, name)

	subchapters := []chapter{}
	if len(configs) > 1 {
		// add subchapters

		links, _, err := GetLinks(base, config.selector, config.limit, config.offset, false)
		if err != nil {
			log.Fatal(err)
		}

		subchapters = make([]chapter, len(links))
		progress := NewProgress(links, name, config.depth)

		if config.delay >= 0 {

			// synchronous mode
			for index, link := range links {
				// and then use it to parse relative URLs
				u, err := base.Parse(link.href)
				if err != nil {
					log.Fatal(err)
				}

				sc := NewChapterFromURL(u.String(), configs[1:], index, progress.UpdateName)
				subchapters[index] = sc
				progress.Increment(index)

				time.Sleep(time.Duration(config.delay) * time.Millisecond)
			}

		} else {
			// asynchronous mode
			var wg sync.WaitGroup

			threads := config.threads
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

					sc := NewChapterFromURL(u.String(), configs[1:], index, progress.UpdateName)
					subchapters[index] = sc
					progress.Increment(index)

					<-semaphore
				}(index, l)
			}
			wg.Wait()
		}
	}

	content := ""
	if config.include {

		// we care about the content only if we include this level
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
				imageTag = strings.ReplaceAll(imageTag, "\n", "")

				content += imageTag
			})

		}
	}

	return chapter{string(body), name, article.Byline, content, subchapters, config}
}

func tableOfContent(url string, config *ScrapeConfig, subConfig *ScrapeConfig) ([]chapter, chapter) {
	base, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	links, home, err := GetLinks(base, config.selector, config.limit, config.offset, config.include)
	if err != nil {
		log.Fatal(err)
	}

	chapters := make([]chapter, len(links))
	// progress := NewProgress(links, "", 0)
	delay := config.delay

	if delay >= 0 {
		// synchronous mode

		for index, link := range links {
			// and then use it to parse relative URLs
			u, err := base.Parse(link.href)
			if err != nil {
				log.Fatal(err)
			}

			sc := NewChapterFromURL(u.String(), []*ScrapeConfig{subConfig}, 0, func(index int, name string) {})
			chapters[index] = sc
			// progress.Increment(index)

			// short sleep for last chapter to let the progress bar update
			if index == len(links)-1 {
				delay = 100
			}

			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

	} else {
		// asynchronous mode
		var wg sync.WaitGroup

		threads := config.threads
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

				sc := NewChapterFromURL(u.String(), []*ScrapeConfig{subConfig}, 0, func(index int, name string) {})
				chapters[index] = sc
				// progress.Increment(index)

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

	home := NewChapterFromURL(url.String(), []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	if include {
		l := NewLink(url.String(), home.Name())
		links = append([]link{l}, links...)
	}

	return links, home, nil
}
