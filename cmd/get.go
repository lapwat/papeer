package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	epub "github.com/bmaupin/go-epub"
	"github.com/spf13/cobra"

	"github.com/lapwat/papeer/book"
)

var recursive, include, images bool
var format, output, selector, name, author string
var limit, offset, delay, threads int

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Scrape URL content",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an URL argument")
		}

		formatEnum := map[string]bool{
			"stdout": true,
			"md":     true,
			"epub":   true,
			"mobi":   true,
		}
		if formatEnum[format] != true {
			return fmt.Errorf("invalid format specified: %s", format)
		}

		// add .mobi to filename if not specified
		if format == "mobi" {
			if len(output) > 0 && strings.HasSuffix(output, ".mobi") == false {
				output = fmt.Sprintf("%s.mobi", output)
			}
		}

		if cmd.Flags().Changed("selector") && recursive == false {
			return errors.New("cannot use selector option if not in recursive mode")
		}

		if cmd.Flags().Changed("include") && recursive == false {
			return errors.New("cannot use include option if not in recursive mode")
		}

		if cmd.Flags().Changed("limit") && recursive == false {
			return errors.New("cannot use limit option if not in recursive mode")
		}

		if cmd.Flags().Changed("offset") && recursive == false {
			return errors.New("cannot use offset option if not in recursive mode")
		}

		if cmd.Flags().Changed("delay") && recursive == false {
			return errors.New("cannot use delay option if not in recursive mode")
		}

		if cmd.Flags().Changed("threads") && recursive == false {
			return errors.New("cannot use threads option if not in recursive mode")
		}

		if cmd.Flags().Changed("delay") && cmd.Flags().Changed("threads") {
			return errors.New("cannot use delay and threads options at the same time")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		b := book.NewBookFromURL(url, selector, name, author, recursive, include, limit, offset, delay, threads)

		if len(output) == 0 {
			// set default output
			output = strings.ReplaceAll(b.Name(), " ", "_")
			output = strings.ReplaceAll(output, "/", "")
			output = fmt.Sprintf("%s.%s", output, format)
		}

		if format == "stdout" {

			for _, c := range b.Chapters() {
				// convert to markdown
				content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
				if err != nil {
					log.Fatal(err)
				}

				text := fmt.Sprintf("%s\n%s\n\n%s\n\n\n", c.Name(), strings.Repeat("=", len(c.Name())), content)

				// write to stdout
				fmt.Println(text)
			}

		}

		if format == "md" {

			// create markdown file
			f, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			for _, c := range b.Chapters() {
				// convert to markdown
				content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
				if err != nil {
					log.Fatal(err)
				}

				text := fmt.Sprintf("%s\n%s\n\n%s\n\n\n", c.Name(), strings.Repeat("=", len(c.Name())), content)

				// write to markdown file
				_, err = f.WriteString(text)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Printf("Markdown saved to \"%s\"\n", output)
		}

		if format == "epub" {
			e := epub.NewEpub(b.Name())
			e.SetAuthor(b.Author())

			for _, c := range b.Chapters() {
				var content string

				if images == false {
					content = c.Content()
				}

				// parse content
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(c.Content()))
				if err != nil {
					log.Fatal(err)
				}

				// retrieve images and download it
				doc.Find("img").Each(func(i int, s *goquery.Selection) {
					src, _ := s.Attr("src")
					imagePath, _ := e.AddImage(src, "")

					if images {
						imageTag, _ := goquery.OuterHtml(s)
						content += imageTag
					}

					content = strings.ReplaceAll(content, src, imagePath)
				})

				html := fmt.Sprintf("<h1>%s</h1>%s", c.Name(), content)
				_, err = e.AddSection(html, c.Name(), "", "")
				if err != nil {
					log.Fatal(err)
				}
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

			for _, c := range b.Chapters() {
				var content string

				if images == false {
					content = c.Content()
				}

				// parse content
				doc, err := goquery.NewDocumentFromReader(strings.NewReader(c.Content()))
				if err != nil {
					log.Fatal(err)
				}

				// retrieve images and download it
				doc.Find("img").Each(func(i int, s *goquery.Selection) {
					src, _ := s.Attr("src")
					imagePath, _ := e.AddImage(src, "")

					if images {
						imageTag, _ := goquery.OuterHtml(s)
						content += imageTag
					}

					content = strings.ReplaceAll(content, src, imagePath)
				})

				html := fmt.Sprintf("<h1>%s</h1>%s", c.Name(), content)
				_, err = e.AddSection(html, c.Name(), "", "")
				if err != nil {
					log.Fatal(err)
				}
			}

			outputEPUB := strings.ReplaceAll(output, ".mobi", ".epub")

			err := e.Write(outputEPUB)
			if err != nil {
				log.Fatal(err)
			}

			exec.Command("kindlegen", outputEPUB).Run()
			// exec command always return status 1 even if it succeed
			// if err != nil {
			// 	log.Fatal(err)
			// }

			fmt.Printf("Ebook saved to \"%s\"\n", output)

			err = os.Remove(outputEPUB)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}
