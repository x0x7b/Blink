package output

import (
	"Blink/types"
	"fmt"
	"strings"
)

func DiffsOutput(results []types.TestResult) {
	for _, res := range results {
		var out strings.Builder
		if len(res.Diffs) == 0 {
			continue
		}
		out.WriteString(types.Cyan + res.Path + types.Reset + " " + res.Payload + "\n")
		for _, d := range res.Diffs {
			if d.Before != "" && d.After != "" {
				out.WriteString(fmt.Sprintf("%v   %-12s %v: %v -> %v\n", types.Magenta, d.Kind, types.Reset, d.Before, d.After))
			} else {
				out.WriteString(fmt.Sprintf("%v   %-12s %v: %v%v\n", types.Magenta, d.Kind, types.Reset, d.Before, d.After))
			}

		}
		fmt.Println(out.String())
	}
}
