package main

import (
	"Blink/core"
	"Blink/scanners"
	"Blink/types"
	"flag"
	"fmt"
	"log"
)

func main() {
	fmt.Println(types.Magenta + "[ Blink v0.5 ]  \n" + types.Reset)
	showBody := flag.Bool("b", false, "Show response body")
	showBody2 := flag.Bool("body", false, "Show response body")
	showBodyLong := flag.Bool("full-body", false, "Show response body")

	method := flag.String("X", "GET", "HTTP method")

	followRedirects := flag.Bool("no-follow", false, "Follow redirects")
	maxRedirects := flag.Int("max-redirects", 5, "Set value to max redirects")

	outputMode := flag.Int("output", 0, "Output mode:\n   0 - default\n   1 - verbose\n   2 - redirect chain")

	timeout := flag.Int("timeout", 5, "Seconds to timeout")

	data := flag.String("data", "", "HTTP Payload")
	urlParam := flag.Bool("url-params", false, "Test URL param for vulns")
	forms := flag.Bool("forms", false, "Test forms on page")

	getSF := flag.Bool("fp", false, "Show server fingerprint")

	flag.Parse()
	var fc types.FlagCondition
	fc.Data = *data
	if *showBody || *showBody2 {
		fc.ShowBody = true
	} else if *showBodyLong {
		fc.ShowFullBody = true
	}
	fc.ShowFp = *getSF
	fc.FollowRedirects = !*followRedirects
	fc.MaxRedirects = *maxRedirects
	fc.Timeout = *timeout
	fc.OutputMode = *outputMode
	fc.TestParam = *urlParam
	fc.TestForms = *forms

	if flag.NArg() < 1 {
		log.Fatal("URL is required")
	}

	url := flag.Arg(0)
	response, redirects, err := core.HttpRequest(*method, url, fc)
	core.ErrorOutput(err)

	if len(redirects) > 0 {
		core.CleanOutput(response, redirects, fc)
	} else {
		core.CleanOutput(response, redirects, fc)
	}
	if fc.TestParam {
		_, results, err := scanners.TesUrlParam(response, fc)
		core.ErrorOutput(err)
		core.CleanOutput(types.BlinkResponse{}, results, fc)
	}
	if fc.TestForms {
		_, results, err := scanners.TestForms(response, fc)
		core.ErrorOutput(err)
		for _, result := range results {
			core.Diffs(result)
		}

	}
}
