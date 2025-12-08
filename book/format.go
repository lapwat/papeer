package book

import (
	"fmt"
	"html"
	"log"
	"os"
	"os/exec"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/PuerkitoBio/goquery"
	epub "github.com/go-shiori/go-epub"
	"github.com/microcosm-cc/bluemonday"
)

func Filename(name string) string {
	filename := name

	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "")

	return filename
}

func ToMarkdownString(c chapter) string {
	markdown := ""

	// chapter content
	if c.config.Include {
		// title
		markdown += fmt.Sprintf("%s\n", c.Name())
		markdown += fmt.Sprintf("%s\n\n", strings.Repeat("=", len(c.Name())))

		// url
		if c.config.PrintURL {
			markdown += fmt.Sprintf("_%s_\n\n", c.URL())
		}

		// convert content to markdown
		content, err := md.ConvertString(c.Content())
		if err != nil {
			log.Fatal(err)
		}
		markdown += fmt.Sprintf("%s\n\n\n", content)
	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		markdown += fmt.Sprintf("%s\n\n\n", ToMarkdownString(sc))
	}

	return markdown
}

func ToMarkdown(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.md", Filename(c.Name()))
	}

	markdown := ToMarkdownString(c)

	// write to file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(markdown)
	if err2 != nil {
		log.Fatal(err2)
	}
	f.Close()

	return filename
}

func ToHtmlString(c chapter) string {
	htmlContent := ""

	// chapter content
	if c.config.Include {
		// title
		htmlContent += fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(c.Name()))

		// url
		if c.config.PrintURL {
			htmlContent += fmt.Sprintf("<p><i>%s</i></p>\n", html.EscapeString(c.URL()))
		}

		// content
		htmlContent += c.Content()
	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		htmlContent += ToHtmlString(sc)
	}

	return htmlContent
}

func ToHtml(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.html", Filename(c.Name()))
	}

	html := fmt.Sprintf("<html><head></head><body>%s</body></html>", ToHtmlString(c))

	// write to file
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, err2 := f.WriteString(html)
	if err2 != nil {
		log.Fatal(err2)
	}
	f.Close()

	return filename
}

func ToEpub(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.epub", Filename(c.Name()))
	}

	// init ebook
	e, err := epub.NewEpub(c.Name())
	if err != nil {
		log.Fatal(err)
	}

	author := c.Author()
	if author == "" {
		author = "Unknown Author"
	}
	e.SetAuthor(author)

	// append table of content
	if len(c.SubChapters()) > 1 {
		html := "<h1>Table of Contents</h1>"

		html += "<ol>"
		for _, sc := range c.SubChapters() {
			html += fmt.Sprintf("<li>%s</li>", sc.Name())
		}
		html += "</ol>"

		_, err := e.AddSection(html, "Table of Contents", "", "")
		if err != nil {
			log.Fatal(err)
		}
	}

	AppendToEpub(e, c)

	err = e.Write(filename)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func AppendToEpub(e *epub.Epub, c chapter) {
	content := ""
	p := bluemonday.UGCPolicy()
	safeHTML := p.Sanitize(c.Content())

	// chapter content
	if c.config.Include {

		if c.config.ImagesOnly == false {
			content = safeHTML
		}

		// parse content
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(safeHTML))
		if err != nil {
			log.Fatal(err)
		}

		// download images and replace src in img tags of content
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			src = strings.Split(src, "?")[0] // remove query part
			imagePath, _ := e.AddImage(src, "")

			// Remove or fix invalid width/height attributes
			s.RemoveAttr("width")
			s.RemoveAttr("height")

			if c.config.ImagesOnly {
				imageTag, _ := goquery.OuterHtml(s)
				content += strings.ReplaceAll(imageTag, src, imagePath)
			} else {
				content = strings.ReplaceAll(content, src, imagePath)
			}
		})

		htmlContent := ""
		// add title only if ImagesOnly = false
		if c.config.ImagesOnly == false {
			htmlContent += fmt.Sprintf("<h1>%s</h1>\n", html.EscapeString(c.Name()))
		}

		// url
		if c.config.PrintURL {
			htmlContent += fmt.Sprintf("<p><i>%s</i></p>\n", html.EscapeString(c.URL()))
		}

		// content
		htmlContent += content

		//  write to epub file
		_, err = e.AddSection(htmlContent, c.Name(), "", "")
		if err != nil {
			log.Fatal(err)
		}

	}

	// subchapters content
	for _, sc := range c.SubChapters() {
		AppendToEpub(e, sc)
	}
}

func ToMobi(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.mobi", Filename(c.Name()))
	} else {

		// add .mobi extension if not specified
		if strings.HasSuffix(filename, ".mobi") == false {
			filename = fmt.Sprintf("%s.mobi", filename)
		}

	}

	filenameEPUB := strings.ReplaceAll(filename, ".mobi", ".epub")
	ToEpub(c, filenameEPUB)

	exec.Command("kindlegen", filenameEPUB).Run()
	// exec command always returns status 1 even if it succeed
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := os.Remove(filenameEPUB)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}
