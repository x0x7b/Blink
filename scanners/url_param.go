package scanners

import (
	"Blink/core"
	"Blink/types"
	"bufio"
	"log"
	"net/url"
	"os"

	"fmt"
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

	for param, _ := range q {
		fmt.Printf("Testing %v\n", param)
		file, ferr := os.Open("wordlists\\urlparam.txt")
		if ferr != nil {
			log.Printf("%s", ferr.Error())
			return response, results, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			new_value := scanner.Text()
			q.Set(param, new_value)
			u.RawQuery = q.Encode()
			newURL := u.String()
			fmt.Printf(types.Magenta+"[SCAN] "+types.Reset+"Testing %s=%s\n"+types.Reset, param, scanner.Text())
			response, _, err = core.HttpRequest(bl.Method, newURL, fc)
			results = append(results, response)
		}
		file.Close()
	}

	return response, results, err
}
