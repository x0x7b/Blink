package main

import (
	"Blink/core"
	"Blink/output"
	"Blink/scanners"
	"Blink/types"
	"flag"
	"fmt"
)

func main() {

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

	ignoreHash := flag.Bool("ignore-hash", false, "Ignore hash diffs")
	ignoreReflection := flag.Bool("ignore-reflection", false, "Ignore reflections")
	wordlist := flag.String("wordlist", "wordlists\\urlparam.txt", "Wordlist for testing")
	top := flag.Int("top", 0, "Show only N results with highest score")

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
	fc.IgnoreHash = *ignoreHash
	fc.IgnoreReflection = *ignoreReflection
	fc.Wordlist = (*wordlist)
	fc.Top = *top

	if flag.NArg() < 1 {
		fmt.Printf(`
    ____  ___       __  
   / __ )/ (_)___  / /__
  / __  / / / __ \/ //_/
 / /_/ / / / / / / ,<   
/_____/_/_/_/ /_/_/|_|  `)
		fmt.Println(types.Magenta + "[ v0.6 ]  \n" + types.Reset)
		fmt.Println(types.Yellow + "[INFO] " + types.Reset + "Use -help to see available options." + types.Reset)
		return
	} else {

		fmt.Println(types.Magenta + "[ Blink v0.6 ]  \n" + types.Reset)
	}

	url := flag.Arg(0)
	response, redirects, err := core.HttpRequest(*method, url, fc)
	output.ErrorOutput(err)

	if len(redirects) > 0 {
		output.CleanOutput(response, redirects, fc)
	} else {
		output.CleanOutput(response, redirects, fc)
	}
	if fc.TestParam {
		results, err := scanners.TesUrlParam(response, fc, output.Report)
		fmt.Printf(types.Yellow + "[WARN] " + types.Reset + "Showing results ONLY with diffs\n")
		if err.Stage != "OK" {
			fmt.Println(output.ErrorOutput(err))
			return
		}
		testresults := core.Diffs(results, fc)
		profile := core.BuildProfile(testresults)
		output.DiffsOutput(testresults, fc)
		output.ProfileOutput(profile)
	}
	if fc.TestForms {
		_, results, err := scanners.TestForms(response, fc, output.Report)
		if err.Stage != "OK" {
			fmt.Println(output.ErrorOutput(err))
			return
		}
		fmt.Printf(types.Yellow + "\n[WARN] " + types.Reset + "Showing results ONLY with diffs\n")
		for _, result := range results {
			results := core.Diffs(result, fc)
			profile := core.BuildProfile(results)
			output.DiffsOutput(results, fc)
			output.ProfileOutput(profile)
		}

	}
	fmt.Print("\033[0m\n")

}
