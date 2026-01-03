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
		if r.RequestData != "" {
			res.Payload = r.RequestData
		} else {
			u, _ := url.Parse(r.URL)
			q := u.Query()
			res.Payload = q.Encode()

		}

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
				if strings.Contains(string(r.Body), value) || strings.Contains(string(r.Body), url.QueryEscape(value)) || strings.Contains(string(r.Body), url.QueryEscape(url.QueryEscape(value))) {
					res.Diffs = append(res.Diffs, diffLine(types.DiffReflect, "", "input reflected", fc))
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
		DiffCookie(baseline.Cookies, r.Cookies, types.FlagCondition{})
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

func DiffCookie(bl []*http.Cookie, res []*http.Cookie, fc types.FlagCondition) []types.Diff {
	var diffs []types.Diff

	blm := CookieMap(bl)
	resm := CookieMap(res)

	for n := range blm {
		_, ok := resm[n]
		if !ok {
			diffs = append(diffs, diffLine(types.DiffCookies, "deleted", "", fc))
			continue
		}
		if blm[n].Raw == resm[n].Raw {
			continue
		}
		bc := blm[n]
		rc := resm[n]

		if bc.Value != rc.Value {
			diffs = append(diffs,
				diffLine(types.DiffCookies, bc.Value, rc.Value, fc))
		}

		if bc.Quoted != rc.Quoted {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					strconv.FormatBool(bc.Quoted),
					strconv.FormatBool(rc.Quoted),
					fc,
				))
		}

		if bc.Path != rc.Path {
			diffs = append(diffs,
				diffLine(types.DiffCookies, bc.Path, rc.Path, fc))
		}

		if bc.Domain != rc.Domain {
			diffs = append(diffs,
				diffLine(types.DiffCookies, bc.Domain, rc.Domain, fc))
		}

		if !bc.Expires.Equal(rc.Expires) {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					bc.Expires.String(),
					rc.Expires.String(),
					fc,
				))
		}

		if bc.MaxAge != rc.MaxAge {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					strconv.Itoa(bc.MaxAge),
					strconv.Itoa(rc.MaxAge),
					fc,
				))
		}

		if bc.Secure != rc.Secure {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					strconv.FormatBool(bc.Secure),
					strconv.FormatBool(rc.Secure),
					fc,
				))
		}

		if bc.HttpOnly != rc.HttpOnly {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					strconv.FormatBool(bc.HttpOnly),
					strconv.FormatBool(rc.HttpOnly),
					fc,
				))
		}

		if bc.SameSite != rc.SameSite {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					sameSiteToString(bc.SameSite),
					sameSiteToString(rc.SameSite),
					fc,
				))
		}

		if bc.Partitioned != rc.Partitioned {
			diffs = append(diffs,
				diffLine(
					types.DiffCookies,
					strconv.FormatBool(bc.Partitioned),
					strconv.FormatBool(rc.Partitioned),
					fc,
				))
		}

	}
	for n := range resm {
		_, ok := blm[n]
		if !ok {
			diffs = append(diffs, diffLine(types.DiffCookies, "added", "", fc))
			continue
		}
	}
	return diffs
}

func CookieMap(cookiem []*http.Cookie) map[string]*http.Cookie {
	m := make(map[string]*http.Cookie, len(cookiem))
	for _, c := range cookiem {
		m[c.Name] = c
	}

	return m
}
func sameSiteToString(s http.SameSite) string {
	switch s {
	case http.SameSiteDefaultMode:
		return "Default"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteNoneMode:
		return "None"
	default:
		return "Unknown"
	}
}
