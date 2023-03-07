package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/elazarl/goproxy"
	readability "github.com/go-shiori/go-readability"
	"github.com/spf13/cobra"
)

type ProxyOptions struct {
	port   int
	output string
}

var proxyOpts *ProxyOptions

func init() {
	proxyOpts = &ProxyOptions{}

	proxyCmd.Flags().IntVarP(&proxyOpts.port, "port", "p", 8080, "Port on which to start the proxy")
	proxyCmd.Flags().StringVarP(&proxyOpts.output, "output", "o", "html", "response format [html, md]")

	rootCmd.AddCommand(proxyCmd)
}

var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Start http proxy",
	Example: "curl --insecure --location --proxy localhost:8080 https://www.eff.org/cyberspace-independence",
	Args: func(cmd *cobra.Command, args []string) error {

		// check provided output is in list
		outputEnum := map[string]bool{
			"html": true,
			"md":   true,
		}
		if outputEnum[proxyOpts.output] != true {
			return fmt.Errorf("invalid output specified: %s", proxyOpts.output)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		proxy := goproxy.NewProxyHttpServer()
		// proxy.Verbose = true

		proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

		proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {

			// extract HTML body
			article, err := readability.FromReader(resp.Body, ctx.Req.URL)
			if err != nil {
				log.Fatal(err)
			}

			content := article.Content

			if proxyOpts.output == "md" {
				// convert content to markdown
				content, err = md.NewConverter("", true, nil).ConvertString(content)
				if err != nil {
					log.Fatal(err)
				}
			}

			stringReader := strings.NewReader(content)
			resp.Body = io.NopCloser(stringReader)

			log.Printf("Serving %s", ctx.Req.URL)

			return resp
		})

		log.Printf("Proxy listening on port %d...", proxyOpts.port)
		log.Printf("Usage: curl --insecure --location --proxy localhost:%d https://www.eff.org/cyberspace-independence", proxyOpts.port)

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", proxyOpts.port), proxy))
	},
}
