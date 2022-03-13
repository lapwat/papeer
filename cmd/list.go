package cmd

import (
	"errors"
	"fmt"
	"log"
	urllib "net/url"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	cobra "github.com/spf13/cobra"

	"github.com/lapwat/papeer/book"
)

type ListOptions struct {
	// url string

	Selector []string
	depth    int
	limit    int
	offset   int
	reverse  bool
	delay    int
	threads  int
	// includeUrl bool
	include     bool
	useLinkName bool
}

var listOpts *ListOptions

func init() {
	listOpts = &ListOptions{}

	// common with get command
	listCmd.Flags().StringSliceVarP(&listOpts.Selector, "selector", "s", []string{}, "table of contents CSS selector")
	listCmd.Flags().IntVarP(&listOpts.depth, "depth", "d", 0, "scraping depth")
	listCmd.Flags().IntVarP(&listOpts.limit, "limit", "l", -1, "limit number of chapters, use with depth/selector")
	listCmd.Flags().IntVarP(&listOpts.offset, "offset", "o", 0, "skip first chapters, use with depth/selector")
	listCmd.Flags().BoolVarP(&listOpts.reverse, "reverse", "r", false, "reverse chapter order")
	listCmd.Flags().IntVarP(&listOpts.delay, "delay", "", -1, "time in milliseconds to wait before downloading next chapter, use with depth/selector")
	listCmd.Flags().IntVarP(&listOpts.threads, "threads", "t", -1, "download concurrency, use with depth/selector")
	listCmd.Flags().BoolVarP(&listOpts.include, "include", "i", false, "include URL as first chapter, use with depth/selector")
	listCmd.Flags().BoolVarP(&listOpts.useLinkName, "use-link-name", "", false, "use link name for chapter title")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:     "list URL",
	Aliases: []string{"ls"},
	Short:   "Print URL table of contents",
	Example: "papeer list https://12factor.net/ -s 'section.concrete>article>h2>a'",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an URL argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(listOpts.Selector) == 0 {
			listOpts.Selector = []string{""}
		}

		base, err := urllib.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		links, path, _, err := book.GetLinks(base, listOpts.Selector[0], listOpts.limit, listOpts.offset, listOpts.reverse, listOpts.include)
		if err != nil {
			log.Fatal(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false

		// format selector path
		pathArray := strings.Split(path, "<")
		// reverse path
		for i, j := 0, len(pathArray)-1; i < j; i, j = i+1, j-1 {
			pathArray[i], pathArray[j] = pathArray[j], pathArray[i]
		}
		pathFormatted := strings.Join(pathArray, ">")

		t.AppendHeader(table.Row{"#", "Name", fmt.Sprintf("Url [%s]", pathFormatted)})

		for index, link := range links {
			u, err := base.Parse(link.Href())
			if err != nil {
				log.Fatal(err)
			}

			t.AppendRow([]interface{}{index + 1, link.Text(), u.String()})
		}

		t.Render()

	},
}
