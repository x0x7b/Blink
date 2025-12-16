package scanners

import (
	"Blink/core"
	"Blink/types"
	"bufio"
	"log"
	"os"

	"fmt"
	"strings"
)

func TesUrlParam(bl types.BlinkResponse, fc types.FlagCondition) (types.BlinkResponse, []types.BlinkResponse, types.BlinkError) {
	fmt.Printf(types.Magenta + "[SCAN] Testing param\n")
	var response types.BlinkResponse
	var redirects []types.BlinkResponse
	var results []types.BlinkResponse
	var err types.BlinkError
	var new_value string
	parts1 := strings.Split(bl.URL, "/")
	parts2 := strings.Split(parts1[len(parts1)-1], "?")
	parts3 := strings.Split(parts2[len(parts2)-1], "=")
	var param, value string
	if len(parts3) == 2 {
		param, value = parts3[0], parts3[1]
	} else {
		return response, redirects, types.BlinkError{Message: "the parameter test flag is enabled, but no parameters were found in the specified url"}
	}

	file, errs := os.Open("wordlists\\urlparam.txt")
	if errs != nil {
		log.Printf(errs.Error())
		return response, results, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		new_value = scanner.Text()
		newURL := strings.Replace(bl.URL, value, new_value, 1)
		fmt.Printf(types.Magenta+"[SCAN] Testing param %s with %s\n"+types.Reset, param, scanner.Text())
		response, _, err = core.HttpRequest(bl.Method, newURL, fc)
		results = append(results, response)
	}

	return response, results, err
}
