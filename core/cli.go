package core

import (
	"flag"
	"fmt"
	"log"
)

func RunCLI() {
	showBody := flag.Bool("b", false, "Show response body")
	showBodyLong := flag.Bool("include-body", false, "Show response body")

	method := flag.String("X", "GET", "HTTP method")

	followRedirects := flag.Bool("no-follow", true, "Follow redirects")
	maxRedirects := flag.Int("max-redirects", 5, "Set value to max redirects")

	timeout := flag.Int("timeout", 5, "Seconds to timeout")

	data := flag.String("data", "", "HTTP Payload")
	flag.Parse()
	var fc FlagCondition
	if *showBody || *showBodyLong {
		fc.ShowBody = true
	}
	fc.FollowRedirects = *followRedirects
	fc.MaxRedirects = *maxRedirects
	fc.Timeout = *timeout

	if flag.NArg() < 1 {
		log.Fatal("URL is required")
	}

	url := flag.Arg(0)
	response, redirects, err := HttpRequest(*method, url, *data, fc)
	if err.Stage != "OK" {
		if err.Stage == "UNKNOWN" {
			fmt.Printf(Red+"[ %v ERROR ] %s\n"+Reset, err.Stage, err.Message)
		} else {
			fmt.Printf(Red+"[ %v ERROR ]\n"+Reset, err.Stage)
		}

		return
	}

	if len(redirects) > 0 {
		for _, redirect := range redirects {
			CleanOutput(redirect, 1, fc)
		}
	}
	CleanOutput(response, 0, fc)
}
