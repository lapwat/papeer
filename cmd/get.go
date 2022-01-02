package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lapwat/papeer/book"
)

var recursive, include, images, quiet bool
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
		b := book.NewBookFromURL(url, selector, name, author, recursive, include, images, quiet, limit, offset, delay, threads)

		fakeConfig := book.NewScrapeConfigFake()
		fakeChapter := book.NewChapter("", b.Name(), b.Author(), "", b.Chapters(), fakeConfig)

		if format == "stdout" {
			// TODO: ToMarkdownString
			markdown := book.ToMarkdown(fakeChapter)
			fmt.Println(markdown)
		}

		if format == "md" {
			// TODO: ToMarkdownFile
			markdown := book.ToMarkdown(fakeChapter)

			if len(output) == 0 {
				filename := book.Filename(fakeChapter.Name())
				output = fmt.Sprintf("%s.md", filename)
			}

			// write to file
			f, err := os.Create(output)
			if err != nil {
				log.Fatal(err)
			}
			_, err2 := f.WriteString(markdown)
			if err2 != nil {
				log.Fatal(err2)
			}
			f.Close()

			fmt.Printf("Markdown saved to \"%s\"\n", output)
		}

		if format == "epub" {
			output = book.ToEpub(fakeChapter, output)
			fmt.Printf("Ebook saved to \"%s\"\n", output)
		}

		if format == "mobi" {
			output = book.ToMobi(fakeChapter, output)
			fmt.Printf("Ebook saved to \"%s\"\n", output)
		}
	},
}
