package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	epub "github.com/bmaupin/go-epub"
	cobra "github.com/spf13/cobra"

	"github.com/papeer/book"
)

var quiet, stdout, recursive, include bool
var format, output, selector string
var delay int

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Scrape URL content",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an URL argument")
		}

		formatEnum := map[string]bool{
			"md":   true,
			"epub": true,
			"mobi": true,
		}
		if formatEnum[format] != true {
			return fmt.Errorf("invalid format specified: %s", format)
		}

		if format == "epub" || format == "mobi" {
			if stdout {
				return errors.New("cannot print EPUB/MOBI file to standard output")
			}
		}

		if format == "mobi" {
			if len(output) > 0 && strings.HasSuffix(output, ".mobi") == false {
				output = fmt.Sprintf("%s.mobi", output)
			}
		}

		if include && recursive == false {
			return errors.New("cannot use include option if not in recursive mode")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		var b book.Book

		if recursive {
			b = book.NewBookFromURL(url, selector, include, delay)
		} else {
			c := book.NewChapterFromURL(url)
			b = book.New(c.Name(), c.Author())
			b.AddChapter(c)
		}

		// if quiet == false {
		// 	metadata := fmt.Sprintf("URL     : %s\nTitle   : %s\nAuthor  : %s\nLength  : %d\nExcerpt : %s\nSiteName: %s\nImage   : %s\nFavicon : %s", url, article.Title, article.Byline, article.Length, article.Excerpt, article.SiteName, article.Image, article.Favicon)
		// 	fmt.Println(metadata)
		// }

		if len(output) == 0 {
			// set default output
			output = strings.ReplaceAll(b.Name(), " ", "_")
			output = strings.ReplaceAll(output, "/", "")
			output = fmt.Sprintf("%s.%s", output, format)
		}

		if format == "md" {
			f, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			for _, c := range b.Chapters() {
				content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
				if err != nil {
					log.Fatal(err)
				}

				if stdout {
					fmt.Println(content)
				} else {

					_, err := f.WriteString(content)
					if err != nil {
						log.Fatal(err)
					}

				}
			}

			fmt.Printf("Markdown saved to \"%s\"\n", output)
		}

		if format == "epub" {
			e := epub.NewEpub(b.Name())
			e.SetAuthor(b.Author())

			for _, c := range b.Chapters() {
				e.AddSection(c.Content(), c.Name(), "", "")
			}

			err := e.Write(output)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Ebook saved to \"%s\"\n", output)
		}

		if format == "mobi" {
			e := epub.NewEpub(b.Name())
			e.SetAuthor(b.Author())

			for _, chapter := range b.Chapters() {
				e.AddSection(chapter.Content(), chapter.Name(), "", "")
			}

			outputEPUB := strings.ReplaceAll(output, ".mobi", ".epub")

			err := e.Write(outputEPUB)
			if err != nil {
				log.Fatal(err)
			}

			exec.Command("kindlegen", outputEPUB).Run()
			// exec command always return status 1 even if it fails
			// if err != nil {
			// 	log.Fatal(err)
			// }

			fmt.Printf("Ebook saved to \"%s\"\n", output)

			err2 := os.Remove(outputEPUB)
			if err2 != nil {
				log.Fatal(err)
			}
		}
	},
}
