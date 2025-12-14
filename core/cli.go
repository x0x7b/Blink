package core

import (
	"flag"
	"fmt"
	"log"
)

func RunCLI() {
	showBody := flag.Bool("b", false, "Show response body")
	showBody2 := flag.Bool("body", false, "Show response body")
	showBodyLong := flag.Bool("full-body", false, "Show response body")

	method := flag.String("X", "GET", "HTTP method")

	followRedirects := flag.Bool("no-follow", false, "Follow redirects")
	maxRedirects := flag.Int("max-redirects", 5, "Set value to max redirects")

	outputMode := flag.Int("output", 0, "Output mode:\n   0 - default\n   1 - verbose\n   2 - redirect chain")

	timeout := flag.Int("timeout", 5, "Seconds to timeout")

	data := flag.String("data", "", "HTTP Payload")
	flag.Parse()
	var fc FlagCondition
	if *showBody || *showBody2 {
		fc.ShowBody = true
	} else if *showBodyLong {
		fc.ShowFullBody = true
	}
	fc.FollowRedirects = !*followRedirects
	fc.MaxRedirects = *maxRedirects
	fc.Timeout = *timeout
	fc.OutputMode = *outputMode

	if flag.NArg() < 1 {
		log.Fatal("URL is required")
	}

	url := flag.Arg(0)
	response, redirects, err := HttpRequest(*method, url, *data, fc)
	if err.Stage != "OK" {
		if err.Stage == "Unknown" {
			fmt.Printf(Red+"[ %v ERROR ] %s\n"+Reset, err.Stage, err.Message)
		} else if err.Stage == "INFO" {
			fmt.Printf(Yellow+"[ %v ] %v"+Reset, err.Stage, err.Message)
		} else {
			fmt.Printf(Red+"[ %v ERROR ]\n"+Reset, err.Stage)
		}

		return
	}

	if len(redirects) > 0 {
		CleanOutput(response, redirects, fc)
	} else {
		CleanOutput(response, redirects, fc)
	}
}
