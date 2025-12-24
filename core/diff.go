package core

import (
	"Blink/types"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func Diffs(bl []types.BlinkResponse, fc types.FlagCondition) []types.TestResult {
	var results []types.TestResult
	if len(bl) == 0 {
		return results
	}
	var baseline = bl[0]
	// var profile types.BehaviorProfile
	// profile.TotalTests = len(bl)
	if baseline.URL == "" {
		return results
	}
	for _, r := range bl[1:] {
		var res types.TestResult
		res.Payload = r.RequestData
		res.Path = r.RawRequest.URL.Path
		if r.URL == "" {
			continue
		}

		if baseline.StatusCode != r.StatusCode {
			res.Diffs = append(res.Diffs, diffLine(types.DiffStatus, strconv.Itoa(baseline.StatusCode), strconv.Itoa(r.StatusCode), fc))
		}

		if baseline.BodyHash != r.BodyHash {
			res.Diffs = append(res.Diffs, diffLine(types.DiffBodyHash, shortHash(baseline.BodyHash), shortHash(r.BodyHash), fc))
		}

		parts := strings.FieldsFunc(r.RequestData, func(g rune) bool {
			return g == '=' || g == '&'
		})

		if len(parts)%2 == 0 {
			for i := 1; i < len(parts); i += 2 {
				value := parts[i]
				if strings.Contains(string(r.Body), value) {
					res.Diffs = append(res.Diffs, diffLine(types.DiffReflect, "", "raw input reflected", fc))
				}

				if strings.Contains(string(r.Body), url.QueryEscape(value)) || strings.Contains(string(r.Body), url.QueryEscape(url.QueryEscape(value))) {
					res.Diffs = append(res.Diffs, diffLine(types.DiffReflect, "", "encoded input reflected", fc))
				}

			}
		}

		headersChanges := diffHeaders(baseline.Headers, r.Headers)
		if len(headersChanges) > 0 {
			res.Diffs = append(res.Diffs, diffLine(types.DiffHeaders, "", strings.Join(headersChanges, ", "), fc))
		}

		if baseline.Timings.FullRtt*2 < r.Timings.FullRtt {
			res.Diffs = append(res.Diffs, diffLine(types.DiffRTT, strconv.FormatInt(int64(baseline.Timings.FullRtt), 10), strconv.FormatInt(int64(r.Timings.FullRtt), 10), fc))
		}

		if len(baseline.Cookies) != len(r.Cookies) {
			res.Diffs = append(res.Diffs, diffLine(types.DiffCookies, strconv.Itoa(len(baseline.Cookies)), strconv.Itoa(len(r.Cookies)), fc))
		}

		results = append(results, res)
	}
	return results
}

func diffHeaders(base, mod http.Header) []string {
	interesting := []string{
		"Content-Type",
		"Location",
		"X-Powered-By",
		"Set-Cookie",
	}
	var changes []string
	for _, h := range interesting {
		if base.Get(h) != mod.Get(h) {
			changes = append(changes, h)
		}
	}
	return changes
}
func diffLine(field types.DiffKind, bfr string, afr string, fc types.FlagCondition) types.Diff {
	var diff types.Diff
	if fc.IgnoreHash {
		if field == types.DiffBodyHash {
			return diff
		}
	}
	if fc.IgnoreReflection {
		if field == types.DiffReflect {
			return diff
		}
	}

	diff.Kind = field
	diff.Before = bfr
	diff.After = afr
	return diff
}

func shortHash(hash string) string {
	if hash != "" {
		return hash[:10]
	} else {
		return "EMPTY_HASH"
	}
}
