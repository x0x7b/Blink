package scanners

import (
	"Blink/core"
	"Blink/types"
	"bufio"
	"log"
	"net/url"
	"os"
)

func TesUrlParam(bl types.BlinkResponse, fc types.FlagCondition) (types.BlinkResponse, []types.BlinkResponse, types.BlinkError) {
	var response types.BlinkResponse
	var redirects []types.BlinkResponse
	var results []types.BlinkResponse
	var err types.BlinkError
	// parts3 := strings.Split(parts2[len(parts2)-1], "=")
	u, _ := url.Parse(bl.URL)
	q := u.Query()
	if len(q) == 0 {
		return response, redirects, types.BlinkError{Message: "the parameter test flag is enabled, but no parameters were found in the specified url"}
	}
	file, ferr := os.Open(fc.Wordlist)
	if ferr != nil {
		log.Printf("%s", ferr.Error())
		return response, results, err
	}

	var payloads []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		payloads = append(payloads, scanner.Text())
	}
	file.Close()
	for param, _ := range q {
		for _, payload := range payloads {
			new_value := payload
			q.Set(param, new_value)
			u.RawQuery = q.Encode()
			newURL := u.String()
			response, _, err = core.HttpRequest(bl.Method, newURL, fc)
			results = append(results, response)
		}
	}

	return response, results, err
}
