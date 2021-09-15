package cmd

import (
	"errors"
	"fmt"
	"log"
	urllib "net/url"
	"strings"

	colly "github.com/gocolly/colly/v2"
	cobra "github.com/spf13/cobra"
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

		if selector == "" {
			selector = "a"
		}

		// visit and count link classes
		classesLinks := map[string][]map[string]string{}
		classesCount := map[string]int{}
		classMax := ""

		c := colly.NewCollector()
		c.OnHTML(selector, func(e *colly.HTMLElement) {
			href := e.Attr("href")
			text := strings.TrimSpace(e.Text)
			class := e.Attr("class")

			if cmd.Flags().Changed("selector") || class != "" && text != "" {
				classesLinks[class] = append(classesLinks[class], map[string]string{
					"href": href,
					"text": text,
				})

				classesCount[class]++

				if classesCount[class] > classesCount[classMax] {
					classMax = class
				}
			}
		})

		c.Visit(base.String())
		for index, link := range classesLinks[classMax] {
			u, err := base.Parse(link["href"])
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Chapter %d: %s %s\n", index+1, link["text"], u.String())
		}

	},
}
