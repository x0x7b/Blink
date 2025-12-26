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
		out.WriteString(fmt.Sprintf("    %-10s: %d (%.0f%%)\n", k, c, profile.Ratios[k]*100))
	}
	fmt.Print(out.String())
}
