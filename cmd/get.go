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

	"github.com/lapwat/papeer/book"
)

var stdout, recursive, include, images bool
var format, output, selector string
var limit, offset, delay int

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

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		b := book.NewBookFromURL(url, selector, recursive, include, images, limit, offset, delay)

		if len(output) == 0 {
			// set default output
			output = strings.ReplaceAll(b.Name(), " ", "_")
			output = strings.ReplaceAll(output, "/", "")
			output = fmt.Sprintf("%s.%s", output, format)
		}

		if format == "md" {
			var f *os.File
			var err error

			if !stdout {
				f, err = os.Create(output)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
			}

			for _, c := range b.Chapters() {
				content, err := md.NewConverter("", true, nil).ConvertString(c.Content())
				if err != nil {
					log.Fatal(err)
				}

				text := fmt.Sprintf("%s\n%s\n\n%s\n\n\n", c.Name(), strings.Repeat("=", len(c.Name())), content)

				if stdout {
					fmt.Println(text)
				} else {
					_, err := f.WriteString(text)
					if err != nil {
						log.Fatal(err)
					}

				}
			}

			if stdout == false {
				fmt.Printf("Markdown saved to \"%s\"\n", output)
			}
		}

		if format == "epub" {
			e := epub.NewEpub(b.Name())
			e.SetAuthor(b.Author())

			for _, c := range b.Chapters() {
				if images {
					e.AddSection(c.Content(), "", "", "")
				} else {
					html := fmt.Sprintf("<h1>%s</h1>%s", c.Name(), c.Content())

					_, err := e.AddSection(html, c.Name(), "", "")
					if err != nil {
						log.Fatal(err)
					}
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
				log.Fatal(err2)
			}
		}
	},
}
