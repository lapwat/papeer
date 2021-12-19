package cmd

import (
	"errors"
	"log"
	urllib "net/url"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	cobra "github.com/spf13/cobra"

	"github.com/lapwat/papeer/book"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "Print table of content",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires an URL argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		base, err := urllib.Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		links, _, err := book.GetLinks(base, selector, limit, offset, include)
		if err != nil {
			log.Fatal(err)
		}

		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.Style().Options.DrawBorder = false
		t.Style().Options.SeparateColumns = false
		t.Style().Options.SeparateHeader = false
		t.AppendHeader(table.Row{"#", "Name", "Url"})

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
