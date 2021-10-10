package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "papeer",
	Short: "Browse the web in the eink era",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "md", "file format [md, epub, mobi]")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "", "", "output file")
	rootCmd.PersistentFlags().StringVarP(&selector, "selector", "s", "", "table of content CSS selector, in resursive mode")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "create one chapter per natigation item")
	rootCmd.PersistentFlags().BoolVarP(&include, "include", "i", false, "include URL as first chapter, in resursive mode")
	rootCmd.PersistentFlags().BoolVarP(&stdout, "stdout", "", false, "print to standard output")
	rootCmd.PersistentFlags().BoolVarP(&images, "images", "", false, "retrieve images only")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", -1, "limit number of chapters, in recursive mode")
	rootCmd.PersistentFlags().IntVarP(&offset, "offset", "o", 0, "skip first chapters, in recursive mode")
	rootCmd.PersistentFlags().IntVarP(&delay, "delay", "d", -1, "time to wait before downloading next chapter, in milliseconds")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", -1, "download concurrency, in recursive mode")

	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
}
