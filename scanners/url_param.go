package scanners

import (
	"Blink/core"
	"Blink/types"

	"fmt"
	"strings"
)

func TesUrlParam(bl types.BlinkResponse, fc types.FlagCondition) (types.BlinkResponse, []types.BlinkResponse, types.BlinkError) {
	var response types.BlinkResponse
	var redirects []types.BlinkResponse
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
	new_value = "'"
	newURL := strings.Replace(bl.URL, value, new_value, 1)
	fmt.Printf("[INFO] Testing param %s\n", param)
	response, redirects, err = core.HttpRequest(bl.Method, newURL, fc)
	return response, redirects, err
}
