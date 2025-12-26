package output

import (
	"Blink/types"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
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
				if d.Kind == types.DiffRTT {
					beforeNS, err := strconv.ParseInt(d.Before, 10, 64)
					if err != nil {
						panic(err)
					}
					afterNS, err := strconv.ParseInt(d.After, 10, 64)
					if err != nil {
						panic(err)
					}

					beforeRTT := time.Duration(beforeNS)
					afterRTT := time.Duration(afterNS)

					beforeMS := beforeRTT.Milliseconds()
					afterMS := afterRTT.Milliseconds()

					out.WriteString(fmt.Sprintf(
						"%v   %-12s %v: %4dms → %4dms\n",
						types.Magenta,
						d.Kind,
						types.Reset,
						beforeMS,
						afterMS,
					))

				} else {
					out.WriteString(fmt.Sprintf("%v   %-12s %v: %v -> %v\n", types.Magenta, d.Kind, types.Reset, d.Before, d.After))
				}

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
