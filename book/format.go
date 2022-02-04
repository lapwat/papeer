package book

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	epub "github.com/bmaupin/go-epub"
)

func Filename(name string) string {
	filename := name

	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ReplaceAll(filename, "/", "")

	return filename
}

func ToMarkdownString(c chapter) string {
	markdown := ""

	if c.config.Include {
		// title
		markdown += fmt.Sprintf("%s\n", c.Name())
		markdown += fmt.Sprintf("%s\n\n", strings.Repeat("=", len(c.Name())))

		// convert content to markdown
		content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
		if err != nil {
			log.Fatal(err)
		}
		markdown += fmt.Sprintf("%s\n\n\n", content)
	}

	for _, sc := range c.SubChapters() {
		// subchapters content
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

func ToEpub(c chapter, filename string) string {
	if len(filename) == 0 {
		filename = fmt.Sprintf("%s.epub", Filename(c.Name()))
	}

	// init ebook
	e := epub.NewEpub(c.Name())
	e.SetAuthor(c.Author())

	AppendToEpub(e, c)

	err := e.Write(filename)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}

func AppendToEpub(e *epub.Epub, c chapter) {
	content := ""

	if c.config.Include {

		if c.config.ImagesOnly == false {
			content = c.Content()
		}

		// parse content
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(c.Content()))
		if err != nil {
			log.Fatal(err)
		}

		// download images and replace src in img tags of content
		doc.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			src = strings.Split(src, "?")[0] // remove query part
			imagePath, _ := e.AddImage(src, "")

			if c.config.ImagesOnly {
				imageTag, _ := goquery.OuterHtml(s)
				content += strings.Replace(imageTag, src, imagePath, 1)
			} else {
				content = strings.Replace(content, src, imagePath, 1)
			}
		})

		html := ""
		// add title only if ImagesOnly = false
		if c.config.ImagesOnly == false {
			html += fmt.Sprintf("<h1>%s</h1>", c.Name())
		}
		html += content

		//  write to epub file
		_, err = e.AddSection(html, c.Name(), "", "")
		if err != nil {
			log.Fatal(err)
		}

	}

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
	// exec command always return status 1 even if it succeed
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := os.Remove(filenameEPUB)
	if err != nil {
		log.Fatal(err)
	}

	return filename
}
