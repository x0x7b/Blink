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
		if res.Payload != "" {
			out.WriteString(fmt.Sprintf("%v %s %v %s %.2f\n", types.Cyan, res.Path, types.Reset, res.Payload, res.Score))
		}

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
