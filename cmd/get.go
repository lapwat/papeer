package cmd

import (
	"errors"
	"fmt"
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
	images bool
	// ImagesOnly bool
	quiet bool

	Selector []string
	depth    int
	limit    int
	offset   int
	delay    int
	threads  int
	// includeUrl bool
	include     bool
	useLinkName bool
}

var getOpts *GetOptions

func init() {
	getOpts = &GetOptions{}

	getCmd.PersistentFlags().StringVarP(&getOpts.name, "name", "n", "", "book name (default: page title)")
	getCmd.PersistentFlags().StringVarP(&getOpts.author, "author", "a", "", "book author")
	getCmd.PersistentFlags().StringVarP(&getOpts.Format, "format", "f", "md", "file format [stdout, md, epub, mobi]")
	getCmd.PersistentFlags().StringVarP(&getOpts.output, "output", "", "", "file name (default: book name)")
	getCmd.PersistentFlags().BoolVarP(&getOpts.images, "images", "", false, "retrieve images only")
	getCmd.PersistentFlags().BoolVarP(&getOpts.quiet, "quiet", "q", false, "hide progress bar")

	// common with list command
	getCmd.Flags().StringSliceVarP(&getOpts.Selector, "selector", "s", []string{}, "table of contents CSS selector")
	getCmd.Flags().IntVarP(&getOpts.depth, "depth", "d", 0, "scraping depth")
	getCmd.Flags().IntVarP(&getOpts.limit, "limit", "l", -1, "limit number of chapters, use with depth/selector")
	getCmd.Flags().IntVarP(&getOpts.offset, "offset", "o", 0, "skip first chapters, use with depth/selector")
	getCmd.Flags().IntVarP(&getOpts.delay, "delay", "", -1, "time in milliseconds to wait before downloading next chapter, use with depth/selector")
	getCmd.Flags().IntVarP(&getOpts.threads, "threads", "t", -1, "download concurrency, use with depth/selector")
	getCmd.Flags().BoolVarP(&getOpts.include, "include", "i", false, "include URL as first chapter, use with depth/selector")
	getCmd.Flags().BoolVarP(&getOpts.useLinkName, "use-link-name", "", false, "use link name for chapter title")

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

		formatEnum := map[string]bool{
			"stdout": true,
			"md":     true,
			"epub":   true,
			"mobi":   true,
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

		if cmd.Flags().Changed("include") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use include option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("limit") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use limit option if depth/selector is not specified")
		}

		if cmd.Flags().Changed("offset") && getOpts.depth == 0 && len(getOpts.Selector) == 0 {
			return errors.New("cannot use offset option if depth/selector is not specified")
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
		url := args[0]

		// fill selector array with empty selectors to match depth
		getOpts.Selector = append(getOpts.Selector, "")
		for len(getOpts.Selector) < getOpts.depth+1 {
			getOpts.Selector = append(getOpts.Selector, "")
		}
		fmt.Println(len(getOpts.Selector))

		// generate config for each level
		configs := make([]*book.ScrapeConfig, len(getOpts.Selector))
		for index, s := range getOpts.Selector {
			config := book.NewScrapeConfig()
			config.Selector = s
			config.Quiet = getOpts.quiet
			config.Limit = getOpts.limit
			config.Offset = getOpts.offset
			config.Delay = getOpts.delay
			config.Threads = getOpts.threads
			config.ImagesOnly = getOpts.images
			config.Include = getOpts.include
			config.UseLinkName = getOpts.useLinkName

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

		c := book.NewChapterFromURL(url, "", configs, 0, func(index int, name string) {})

		if getOpts.Format == "stdout" {
			markdown := book.ToMarkdownString(c)
			fmt.Println(markdown)
		}

		if getOpts.Format == "md" {
			filename := book.ToMarkdown(c, getOpts.output)
			fmt.Printf("Markdown saved to \"%s\"\n", filename)
		}

		if getOpts.Format == "epub" {
			filename := book.ToEpub(c, getOpts.output)
			fmt.Printf("Ebook saved to \"%s\"\n", filename)
		}

		if getOpts.Format == "mobi" {
			filename := book.ToMobi(c, getOpts.output)
			fmt.Printf("Ebook saved to \"%s\"\n", filename)
		}
	},
}
