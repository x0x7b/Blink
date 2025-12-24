package core

import "Blink/types"

func BuildProfile(results []types.TestResult) types.BehaviorProfile {
	var profile types.BehaviorProfile
	profile.TotalTests = len(results)
	counts := make(map[types.DiffKind]int)
	for _, r := range results {
		for _, d := range r.Diffs {
			counts[d.Kind] += 1
		}
	}
	profile.Counts = counts

	return profile
}
