package output

import (
	"Blink/types"
	"fmt"
	"strings"
)

func redirectChainOutput(redirects []types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	var stringWidth int
	out.WriteString(types.Cyan + "Redirect chain:\n" + types.Reset)
	if len(redirects) > 0 {
		var maxLenRequest int
		for i, req := range redirects {
			var outForLen strings.Builder
			outForLen.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outForLen.WriteString(fmt.Sprintf("%v ", req.StatusCode))
			outForLen.WriteString("[ " + req.Method + " " + req.URL + " ] ")
			outForLen.WriteString(fmt.Sprintf("(%v)", req.Timings.FullRtt))
			if len(outForLen.String()) > maxLenRequest {
				maxLenRequest = len(outForLen.String())
			}
			stringWidth = maxLenRequest + 1
		}

		for i, redirect := range redirects {
			var outString strings.Builder
			var outNoColors strings.Builder
			outString.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outString.WriteString(ColorStatus(redirect.StatusCode) + fmt.Sprintf("%v ", redirect.StatusCode) + types.Reset)
			outString.WriteString(types.Blue + "[ " + types.Reset + types.Cyan + redirect.Method + types.Reset + " " + redirect.URL + types.Blue + " ] " + types.Reset)
			outString.WriteString(fmt.Sprintf("(%v)", redirect.Timings.FullRtt))

			outNoColors.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outNoColors.WriteString(fmt.Sprintf("%v ", redirect.StatusCode))
			outNoColors.WriteString("[ " + redirect.Method + " " + redirect.URL + " ] ")
			outNoColors.WriteString(fmt.Sprintf("(%v)", redirect.Timings.FullRtt))

			spaces := stringWidth - len(outNoColors.String())
			if spaces <= 0 {
				spaces = 1
			}
			outString.WriteString(fmt.Sprintf("%v> %v\n", strings.Repeat(" ", spaces), redirect.Headers.Get("Location")))
			out.WriteString(outString.String())
		}
	} else {
		fmt.Println("No redirects")
	}
	fmt.Print(out.String())

}
