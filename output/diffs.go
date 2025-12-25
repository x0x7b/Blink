package output

import (
	"Blink/types"
	"fmt"
	"sort"
	"strings"
)

func DiffsOutput(results []types.TestResult, fc types.FlagCondition) {

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score // DESC
	})
	if fc.Top == 0 {
		fc.Top = len(results)
	}
	for i, res := range results {
		var out strings.Builder
		if len(res.Diffs) == 0 {
			continue
		}
		if res.Payload != "" {
			out.WriteString(fmt.Sprintf("%v %s %v %.2f %v %v\n", types.Cyan, res.Path, colorScore(res.Score), res.Score, types.Reset, res.Payload))
		}

		for _, d := range res.Diffs {
			if d.Before != "" && d.After != "" {
				out.WriteString(fmt.Sprintf("%v   %-12s %v: %v -> %v\n", types.Magenta, d.Kind, types.Reset, d.Before, d.After))
			} else {
				out.WriteString(fmt.Sprintf("%v   %-12s %v: %v%v\n", types.Magenta, d.Kind, types.Reset, d.Before, d.After))
			}

		}
		if i+1 <= fc.Top {
			fmt.Println(out.String())
		} else {
			break
		}

	}
}
