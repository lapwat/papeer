package book

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	urllib "net/url"
	"strings"
	"sync"
	"time"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	colly "github.com/gocolly/colly/v2"
	"github.com/mmcdole/gofeed"
)

type ScrapeConfig struct {
	Depth       int
	Selector    string
	Quiet       bool
	Limit       int
	Offset      int
	Reverse     bool
	Delay       int
	Threads     int
	Include     bool
	ImagesOnly  bool
	UseLinkName bool
	PrintURL    bool
	Browser     bool
}

func NewScrapeConfig() *ScrapeConfig {
	return &ScrapeConfig{0, "", false, -1, 0, false, -1, -1, true, false, false, false, false}
}

func NewScrapeConfigQuiet() *ScrapeConfig {
	return &ScrapeConfig{0, "", true, -1, 0, false, -1, -1, true, false, false, false, false}
}

func NewScrapeConfigNoInclude() *ScrapeConfig {
	return &ScrapeConfig{0, "", false, -1, 0, false, -1, -1, false, false, false, false, false}
}

func NewScrapeConfigs(selectors []string) []*ScrapeConfig {
	configs := []*ScrapeConfig{}

	for _, s := range selectors {
		config := NewScrapeConfig()
		config.Selector = s

		configs = append(configs, config)
	}

	return configs
}

func NewScrapeConfigsAjin() []*ScrapeConfig {
	config0 := NewScrapeConfig()
	config0.Depth = 0
	config0.Selector = ".dt>a"
	config0.Limit = 3
	config0.Offset = 0
	config0.Delay = 5000
	config0.Include = false

	config1 := NewScrapeConfig()
	config1.Depth = 1
	config1.Selector = ".nav_apb>a"
	config1.Limit = 3
	config1.Offset = 1
	config1.Delay = 5000
	config1.Include = false

	config2 := NewScrapeConfig()
	config2.Depth = 2
	config2.ImagesOnly = true

	return []*ScrapeConfig{config0, config1, config2}
}

func NewScrapeConfigsWikipedia() []*ScrapeConfig {
	config0 := NewScrapeConfig()
	config0.Depth = 0
	config0.Threads = -1
	config0.Include = true

	config1 := NewScrapeConfig()
	config1.Depth = 1
	config1.Include = true

	return []*ScrapeConfig{config0, config1}
}

func NewScrapeConfigFake() *ScrapeConfig {
	config := NewScrapeConfig()
	config.Include = false

	return config
}

var (
	browserCtx  context.Context
	browserOnce sync.Once
)

func getBrowserContext() context.Context {
	browserOnce.Do(func() {
		allocCtx, _ := chromedp.NewExecAllocator(context.Background(),
			append(
				chromedp.DefaultExecAllocatorOptions[:],
				chromedp.NoSandbox,
				chromedp.UserAgent("papeer"),
			)...,
		)
		browserCtx, _ = chromedp.NewContext(allocCtx)
		chromedp.Run(browserCtx)
	})
	return browserCtx
}

func fetchHTMLStd(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "papeer")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func fetchHTMLWithBrowser(url string) (io.ReadCloser, error) {
	ctx, cancel := chromedp.NewContext(getBrowserContext())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return nil, fmt.Errorf("headless browser: %w", err)
	}

	return io.NopCloser(strings.NewReader(htmlContent)), nil
}

func NewChapterFromURL(url, linkName string, configs []*ScrapeConfig, index int, updateProgressBarName func(index int, name string)) chapter {
	config := configs[0]

	baseUrl, err := urllib.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	// fetch content
	var response io.ReadCloser
	if config.Browser {
		response, err = fetchHTMLWithBrowser(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Close()
	} else {
		response, err = fetchHTMLStd(url)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Close()
	}

	// duplicate response stream
	readabilityReader := &bytes.Buffer{}
	bodyReader := io.TeeReader(response, readabilityReader)

	// extract HTML body
	body, err := io.ReadAll(bodyReader)

	// extract article content and metadata
	article, err := readability.FromReader(readabilityReader, baseUrl)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", url, err)
	}

	name := linkName
	if config.UseLinkName == false {
		name = article.Title()

		// notify progressbar with new name
		updateProgressBarName(index, name)
	}

	var subchapters []chapter
	if len(configs) > 1 {

		// retrieve links on page
		links, _, _, err := GetLinks(baseUrl, config.Selector, config.Limit, config.Offset, config.Reverse, false, config.Browser)
		if err != nil {
			log.Fatal(err)
		}

		// init progress bar
		var p progress
		if config.Quiet == false {
			p = NewProgress(links, name, config.Depth)
		}

		// init chapters list
		subchapters = make([]chapter, len(links))

		if config.Delay >= 0 {

			// synchronous mode
			for index, link := range links {
				// and then use it to parse relative URLs
				u, err := baseUrl.Parse(link.Href)
				if err != nil {
					log.Fatal(err)
				}

				sc := NewChapterFromURL(u.String(), link.Text, configs[1:], index, p.UpdateName)
				subchapters[index] = sc
				if config.Quiet == false {
					p.Increment(index)
				}

				time.Sleep(time.Duration(config.Delay) * time.Millisecond)
			}

		} else {
			// asynchronous mode
			var wg sync.WaitGroup

			threads := config.Threads
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
					u, err := baseUrl.Parse(l.Href)
					if err != nil {
						log.Fatal(err)
					}

					sc := NewChapterFromURL(u.String(), l.Text, configs[1:], index, p.UpdateName)
					subchapters[index] = sc

					if config.Quiet == false {
						p.Increment(index)
					}

					<-semaphore
				}(index, l)
			}
			wg.Wait()
		}
	}

	content := ""
	if config.Include {

		// we care about the content only if:
		// - we include this level
		// - we use the page name

		var buffer bytes.Buffer
		err := article.RenderHTML(&buffer)
		if err != nil {
			log.Fatal(err)
		}

		// parse HTML
		doc, err := goquery.NewDocumentFromReader(&buffer)
		if err != nil {
			log.Fatal(err)
		}

		// handle lazy images
		doc.Find("img").Each(func(i int, source *goquery.Selection) {
			src, exists := source.Attr("data-lazy-src")
			if exists {
				source.SetAttr("src", src)
			}
		})
		doc.Find("source").Remove()

		// extract images
		if config.ImagesOnly {

			// append every image to content
			content = ""
			doc.Find("img").Each(func(i int, s *goquery.Selection) {
				imageTag, _ := goquery.OuterHtml(s)
				content += imageTag
			})

		} else {

			content, err = doc.Find("[id*=readability-page]").Html()
			if err != nil {
				log.Fatal(err)
			}

		}

	}

	return chapter{url, string(body), name, article.Byline(), content, subchapters, config}
}

func GetPath(elm *goquery.Selection) string {
	path := []string{}

	for {
		selector := strings.ToLower(goquery.NodeName(elm))
		if len(selector) == 0 {
			break
		}

		path = append(path, selector)
		elm = elm.Parent()
	}

	join := strings.Join(path, "<")
	return join
}

func GetLinks(url *urllib.URL, selector string, limit, offset int, reverse, include, Browser bool) ([]link, string, chapter, error) {
	var links []link
	var pathMax string

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url.String())

	if err == nil {
		// RSS feed

		for _, item := range feed.Items {
			u, err := url.Parse(item.Link)
			if err != nil {
				return []link{}, "", chapter{}, err
			}

			links = append(links, NewLink(u.String(), item.Title, item.PublishedParsed))
		}

		pathMax = "RSS"
	} else {
		// HTML website

		selectorSet := true
		if len(selector) == 0 {
			selector = "a"
			selectorSet = false
		}

		pathLinks := map[string][]link{}
		pathCount := map[string]int{}
		pathMax = ""

		linkHandler := func(s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			path := GetPath(s)
			key := path

			hrefStr := s.AttrOr("href", "")
			u, err := url.Parse(hrefStr)
			if err != nil {
				log.Fatal(err)
			}
			href := u.String()

			if selectorSet {
				// if selector is set, we use the selector specified by the user
				key = selector
				pathLinks[key] = append(pathLinks[key], NewLink(href, text, &time.Time{}))
				pathCount[key] += 1
				pathMax = key
			} else {
				// if selector is not set, we compute the selector ourselves
				class := s.AttrOr("class", "")

				// include the element class to make sure we have the same exact path for every link in the table of content
				key = fmt.Sprintf("%s.%s", path, class)

				// we count this key if the link text is not empty
				if text != "" {
					pathLinks[key] = append(pathLinks[key], NewLink(href, text, &time.Time{}))
					pathCount[key] += len(text)

					if pathCount[key] > pathCount[pathMax] {
						pathMax = key
					}
				}
			}
		}

		if Browser {
			htmlReader, err := fetchHTMLWithBrowser(url.String())
			if err != nil {
				return []link{}, "", chapter{}, err
			}
			defer htmlReader.Close()

			doc, err := goquery.NewDocumentFromReader(htmlReader)
			if err != nil {
				return []link{}, "", chapter{}, err
			}

			doc.Find(selector).Each(func(i int, s *goquery.Selection) {
				linkHandler(s)
			})
		} else {
			c := colly.NewCollector()
			c.OnHTML(selector, func(e *colly.HTMLElement) {
				linkHandler(e.DOM)
			})
			c.Visit(url.String())
		}

		links = pathLinks[pathMax]
	}

	if len(links) == 0 {
		return []link{}, pathMax, chapter{}, fmt.Errorf("no link found for selector: %s", selector)
	}

	end := len(links)
	if limit != -1 {
		end = int(math.Min(float64(limit+offset), float64(len(links))))
	}

	links = links[offset:end]

	home := NewChapterFromURL(url.String(), "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	// include home page
	if include {
		l := NewLink(url.String(), home.Name(), &time.Time{})
		links = append([]link{l}, links...)
	}

	// reverse links
	if reverse {
		for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
			links[i], links[j] = links[j], links[i]
		}
	}

	return links, pathMax, home, nil
}
