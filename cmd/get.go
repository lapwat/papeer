package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lapwat/papeer/book"
)

type GetOptions struct {
	// url string

	name   string
	author string
	Format string
	output string
	stdout bool
	images bool
	quiet  bool

	Selector    []string
	depth       int
	limit       int
	offset      int
	reverse     bool
	delay       int
	threads     int
	include     bool
	useLinkName bool
	printURL    bool
	useHeadless bool // Use headless browser for JavaScript-rendered content
}

var getOpts *GetOptions

func init() {
	getOpts = &GetOptions{}

	getCmd.Flags().StringVarP(&getOpts.name, "name", "n", "", "book name (default: page title)")
	getCmd.Flags().StringVarP(&getOpts.author, "author", "a", "", "book author")
	getCmd.Flags().StringVarP(&getOpts.Format, "format", "f", "md", "file format [md, json, html, epub, mobi]")
	getCmd.Flags().StringVarP(&getOpts.output, "output", "", "", "file name (default: book name)")
	getCmd.Flags().BoolVarP(&getOpts.stdout, "stdout", "", false, "print to standard output")
	getCmd.Flags().BoolVarP(&getOpts.images, "images", "", false, "retrieve images only")
	getCmd.Flags().BoolVarP(&getOpts.printURL, "print-url", "", false, "print url after chapter title")
	getCmd.Flags().BoolVarP(&getOpts.quiet, "quiet", "q", false, "hide progress bar")

	// common with list command
	getCmd.Flags().StringSliceVarP(&getOpts.Selector, "selector", "s", []string{}, "table of contents CSS selector")
	getCmd.Flags().IntVarP(&getOpts.depth, "depth", "d", 0, "scraping depth")
	getCmd.Flags().IntVarP(&getOpts.limit, "limit", "l", -1, "limit number of chapters, use with depth/selector")
	getCmd.Flags().IntVarP(&getOpts.offset, "offset", "o", 0, "skip first chapters, use with depth/selector")
	getCmd.Flags().BoolVarP(&getOpts.reverse, "reverse", "r", false, "reverse chapter order")
	getCmd.Flags().IntVarP(&getOpts.delay, "delay", "", -1, "time in milliseconds to wait before downloading next chapter, use with depth/selector")
	getCmd.Flags().IntVarP(&getOpts.threads, "threads", "t", -1, "download concurrency, use with depth/selector")
	getCmd.Flags().BoolVarP(&getOpts.include, "include", "i", false, "include URL as first chapter, use with depth/selector")
	getCmd.Flags().BoolVarP(&getOpts.useLinkName, "use-link-name", "", false, "use link name for chapter title")
	getCmd.Flags().BoolVarP(&getOpts.useHeadless, "headless", "", false, "use headless browser for JavaScript-rendered content")

	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:     "get URL",
	Short:   "Scrape URL content",
	Example: "papeer get https://www.eff.org/cyberspace-independence",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an URL argument")
		}

		// check provided format is in list
		formatEnum := map[string]bool{
			"md":   true,
			"json": true,
			"html": true,
			"epub": true,
			"mobi": true,
		}
		if formatEnum[getOpts.Format] != true {
			return fmt.Errorf("invalid format specified: %s", getOpts.Format)
		}

		// add .mobi to filename if not specified
		if getOpts.Format == "mobi" {
			if len(getOpts.output) > 0 && strings.HasSuffix(getOpts.output, ".mobi") == false {
				getOpts.output = fmt.Sprintf("%s.mobi", getOpts.output)
			}
		}

		// increase depth to match limit
		if cmd.Flags().Changed("limit") && getOpts.depth == 0 {
			getOpts.depth = 1
		}

		// fill selector array with empty selectors to match depth
		getOpts.Selector = append(getOpts.Selector, "")
		for len(getOpts.Selector) < getOpts.depth+1 {
			getOpts.Selector = append(getOpts.Selector, "")
		}

		if cmd.Flags().Changed("include") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use include option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("offset") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use offset option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("reverse") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use reverse option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("delay") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use delay option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("threads") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use threads option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("use-link-name") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use use-link-name option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("delay") && cmd.Flags().Changed("threads") {
			return errors.New("cannot use delay and threads options at the same time")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// generate config for each level
		configs := make([]*book.ScrapeConfig, len(getOpts.Selector))
		for index, s := range getOpts.Selector {
			config := book.NewScrapeConfig()
			config.Selector = s
			config.Quiet = getOpts.quiet
			config.Limit = getOpts.limit
			config.Offset = getOpts.offset
			config.Reverse = getOpts.reverse
			config.Delay = getOpts.delay
			config.Threads = getOpts.threads
			config.ImagesOnly = getOpts.images
			config.Include = getOpts.include
			config.UseLinkName = getOpts.useLinkName
			config.PrintURL = getOpts.printURL
			config.UseHeadless = getOpts.useHeadless

			// do not use link name for root level as there is not parent link
			if index == 0 {
				config.UseLinkName = false
			}

			// always include last level by default
			if index == len(getOpts.Selector)-1 {
				config.Include = true
			}

			configs[index] = config
		}

		// dummy root chapter to contain all subchapters
		c := book.NewEmptyChapter()
		for _, u := range args {
			newChapter := book.NewChapterFromURL(u, "", configs, 0, func(index int, name string) {})
			c.AddSubChapter(newChapter)
		}
		c.SetName(c.SubChapters()[0].Name())

		if getOpts.Format == "md" {
			filename := book.ToMarkdown(c, getOpts.output)

			if getOpts.stdout {
				bytesRead, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				os.Remove(filename)

				fmt.Println(string(bytesRead))
			} else {
				fmt.Printf("Markdown saved to \"%s\"\n", filename)
			}
		}

		if getOpts.Format == "json" {
			filename := book.ToMarkdown(c, getOpts.output)

			bytesRead, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}
			os.Remove(filename)

			book := make(map[string]interface{})
			book["name"] = c.Name()
			book["content"] = string(bytesRead)

			bookJson, err := json.Marshal(book)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(bookJson))
		}

		if getOpts.Format == "html" {
			filename := book.ToHtml(c, getOpts.output)

			if getOpts.stdout {
				bytesRead, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				os.Remove(filename)

				fmt.Println(string(bytesRead))
			} else {
				fmt.Printf("Html saved to \"%s\"\n", filename)
			}
		}

		if getOpts.Format == "epub" {
			filename := book.ToEpub(c, getOpts.output)

			if getOpts.stdout {
				bytesRead, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				os.Remove(filename)

				fmt.Println(string(bytesRead))
			} else {
				fmt.Printf("Ebook saved to \"%s\"\n", filename)
			}
		}

		if getOpts.Format == "mobi" {
			filename := book.ToMobi(c, getOpts.output)

			if getOpts.stdout {
				bytesRead, err := ioutil.ReadFile(filename)
				if err != nil {
					log.Fatal(err)
				}
				os.Remove(filename)

				fmt.Println(string(bytesRead))
			} else {
				fmt.Printf("Ebook saved to \"%s\"\n", filename)
			}
		}
	},
}
