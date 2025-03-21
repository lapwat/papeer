package book

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

// GetLinksWithHeadlessBrowser gets links from a JavaScript-rendered page
func GetLinksWithHeadlessBrowser(url *url.URL, selector string, limit, offset int, reverse, include bool) ([]link, string, chapter, error) {
	// Create a new context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Navigate to the page and wait for it to load
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url.String()),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		// Wait a bit for JavaScript to execute
		chromedp.Sleep(2*time.Second),
		// Get the HTML content
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return []link{}, "", chapter{}, fmt.Errorf("failed to navigate to %s: %v", url.String(), err)
	}

	// Parse the HTML with goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return []link{}, "", chapter{}, fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Use the same logic as GetLinks to extract links
	var links []link
	pathLinks := map[string][]link{}
	pathCount := map[string]int{}
	pathMax := ""

	selectorSet := true
	if len(selector) == 0 {
		selector = "a"
		selectorSet = false
	}

	// Find all matching elements
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		path := GetPath(s)
		key := path

		href, exists := s.Attr("href")
		if !exists {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			log.Printf("Error parsing URL %s: %v", href, err)
			return
		}
		href = u.String()

		if selectorSet {
			// if selector is set, we use the selector specified by the user
			key = selector
			pathLinks[key] = append(pathLinks[key], NewLink(href, text, &time.Time{}))
			pathCount[key] += 1
			pathMax = key
		} else {
			// if selector is not set, we compute the selector ourselves
			class, _ := s.Attr("class")
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
	})

	links = pathLinks[pathMax]

	if len(links) == 0 {
		return []link{}, pathMax, chapter{}, fmt.Errorf("no link found for selector: %s", selector)
	}

	// Apply limit and offset
	end := len(links)
	if limit != -1 && offset+limit < end {
		end = offset + limit
	}

	if offset < len(links) {
		links = links[offset:end]
	} else {
		links = []link{}
	}

	// Create home chapter
	home := NewChapterFromURL(url.String(), "", []*ScrapeConfig{NewScrapeConfig()}, 0, func(index int, name string) {})

	// Include home page
	if include {
		l := NewLink(url.String(), home.Name(), &time.Time{})
		links = append([]link{l}, links...)
	}

	// Reverse links
	if reverse {
		for i, j := 0, len(links)-1; i < j; i, j = i+1, j-1 {
			links[i], links[j] = links[j], links[i]
		}
	}

	return links, pathMax, home, nil
}
