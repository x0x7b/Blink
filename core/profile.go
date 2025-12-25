package core

import (
	"Blink/types"
)

func BuildProfile(results []types.TestResult) types.BehaviorProfile {
	var profile types.BehaviorProfile
	profile.TotalTests = len(results)
	if profile.TotalTests == 0 {
		return profile
	}
	counts := make(map[types.DiffKind]int)

	for _, r := range results {
		seen := make(map[types.DiffKind]bool)

		for _, d := range r.Diffs {
			seen[d.Kind] = true
		}

		for k := range seen {
			counts[k]++
		}
	}

	ratios := make(map[types.DiffKind]float64)
	for k, _ := range counts {
		ratios[k] = float64(counts[k]) / float64(profile.TotalTests)
	}
	profile.Ratios = ratios
	profile.Counts = counts

	for i := range results {
		var score float64
		for _, d := range (results)[i].Diffs {
			score += types.BaseWeight[d.Kind] * (1 - ratios[d.Kind])
		}
		(results)[i].Score = score
	}

	return profile
}
