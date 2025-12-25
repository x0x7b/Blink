package output

import (
	"Blink/types"
	"fmt"
	"strconv"
	"strings"
)

func ProfileOutput(profile types.BehaviorProfile) {
	// j, _ := json.MarshalIndent(profile, "", "    ")
	// fmt.Println(string(j))
	var out strings.Builder
	out.WriteString(types.Cyan + "Behavior profile:\n" + types.Reset)
	out.WriteString("  Total tests: " + strconv.Itoa(profile.TotalTests) + "\n")
	for k, c := range profile.Counts {
		out.WriteString(fmt.Sprintf("    %-10s: %d (%.2f)\n", k, c, profile.Ratios[k]))
	}
	fmt.Print(out.String())
}
