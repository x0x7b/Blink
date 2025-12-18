package core

import (
	"Blink/types"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func ColorStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return types.Green
	case code >= 300 && code < 400:
		return types.Blue
	case code >= 400 && code < 500:
		return types.Yellow
	default:
		return types.Red
	}
}

func colorTime(timing time.Duration) string {
	switch {
	case timing < 150*time.Millisecond:
		return types.Green
	case timing >= 150*time.Microsecond:
		return types.Yellow
	case timing >= 300*time.Millisecond:
		return types.Red
	default:
		return types.Red
	}

}

func CleanOutput(bl types.BlinkResponse, rc []types.BlinkResponse, fc types.FlagCondition) {
	if len(rc) > 0 && bl.URL != "" {
		redirectChainOutput(rc, fc)
	}
	if len(rc) > 0 && bl.URL == "" {
		Diffs(rc, fc)
	}
	if bl.URL == "" {
		return
	}
	switch {
	case fc.OutputMode == 0:
		defaultOutput(bl, fc)
	case fc.OutputMode == 1:
		verboseOutput(bl, fc)
	}

}

func defaultOutput(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	if fc.OutputMode != 0 {
		out.WriteString(types.Cyan + types.Bold + "Final response: \n" + types.Reset)
	}
	out.WriteString(ColorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + types.Reset)
	out.WriteString(types.Blue + "[ " + types.Reset + types.Cyan + bl.Method + types.Reset + " " + bl.URL + types.Blue + " ] " + types.Reset)
	out.WriteString(fmt.Sprintf("(%v)\n", bl.Timings.FullRtt))
	out.WriteString(bodyOutput(bl, fc))
	serverFingerprint(bl, fc)
	fmt.Print(out.String())

}

func verboseOutput(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	if fc.OutputMode != 0 {
		out.WriteString(types.Cyan + types.Bold + "Final response: \n" + types.Reset)
	}

	out.WriteString(fmt.Sprintf(
		types.Bold+"method: "+types.Reset+"%s\n"+
			types.Bold+"url:    "+types.Reset+"%s\n"+
			types.Bold+"status: "+types.Reset+ColorStatus(bl.StatusCode)+"%d (%s)"+types.Reset+"\n"+
			types.Bold+"proto:  "+types.Reset+"%s (%d.%d)\n"+
			types.Bold+"rtt:    "+types.Reset+"%s\n"+
			types.Bold+"  dns:    "+types.Reset+colorTime(bl.Timings.DnsDuration)+"%s\n"+types.Reset+
			types.Bold+"  tcp:    "+types.Reset+colorTime(bl.Timings.TcpDuration)+"%s\n"+types.Reset+
			types.Bold+"  tls:    "+types.Reset+colorTime(bl.Timings.TlsDuration)+"%s\n"+types.Reset+
			types.Bold+"  ttfb:   "+types.Reset+colorTime(bl.Timings.Ttfb)+"%s\n"+types.Reset+
			types.Bold+"  length: "+types.Reset+"%d\n",
		bl.Method,
		bl.URL,
		bl.StatusCode, bl.Status,
		bl.Proto, bl.ProtoMajor, bl.ProtoMinor,
		bl.Timings.FullRtt,
		bl.Timings.DnsDuration,
		bl.Timings.TcpDuration,
		bl.Timings.TlsDuration,
		bl.Timings.Ttfb,
		bl.ContentLength,
	))

	out.WriteString(types.Cyan + "TLS:" + types.Reset + "\n")
	out.WriteString(fmt.Sprintf("   alpn: %v\n", bl.ALPN))
	out.WriteString(fmt.Sprintf("   Version: %v\n", bl.TLSVersion))
	out.WriteString(fmt.Sprintf("   Cipher:  %v\n", bl.CipherSuite))
	out.WriteString(fmt.Sprintf("   Issuer:  %v\n", bl.CertIssuer))
	out.WriteString(fmt.Sprintf("   Expires: %v\n", bl.CertExpires))

	out.WriteString(types.Cyan + "headers:" + types.Reset + "\n")
	for k, v := range bl.Headers {
		out.WriteString(fmt.Sprintf("  "+types.Cyan+"%s:"+types.Reset+" %s\n", k, v))
	}

	out.WriteString(bodyOutput(bl, fc))
	fmt.Print(out.String())
	serverFingerprint(bl, fc)
}

func redirectChainOutput(redirects []types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	var stringWidth int
	out.WriteString(types.Cyan + "Redirect chain:\n" + types.Reset)
	if len(redirects) > 0 {
		var maxLenRequest int
		for i, req := range redirects {
			var outForLen strings.Builder
			outForLen.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outForLen.WriteString(fmt.Sprintf("%v ", req.StatusCode))
			outForLen.WriteString("[ " + req.Method + " " + req.URL + " ] ")
			outForLen.WriteString(fmt.Sprintf("(%v)", req.Timings.FullRtt))
			if len(outForLen.String()) > maxLenRequest {
				maxLenRequest = len(outForLen.String())
			}
			stringWidth = maxLenRequest + 1
		}

		for i, redirect := range redirects {
			var outString strings.Builder
			var outNoColors strings.Builder
			outString.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outString.WriteString(ColorStatus(redirect.StatusCode) + fmt.Sprintf("%v ", redirect.StatusCode) + types.Reset)
			outString.WriteString(types.Blue + "[ " + types.Reset + types.Cyan + redirect.Method + types.Reset + " " + redirect.URL + types.Blue + " ] " + types.Reset)
			outString.WriteString(fmt.Sprintf("(%v)", redirect.Timings.FullRtt))

			outNoColors.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outNoColors.WriteString(fmt.Sprintf("%v ", redirect.StatusCode))
			outNoColors.WriteString("[ " + redirect.Method + " " + redirect.URL + " ] ")
			outNoColors.WriteString(fmt.Sprintf("(%v)", redirect.Timings.FullRtt))

			spaces := stringWidth - len(outNoColors.String())
			if spaces <= 0 {
				spaces = 1
			}
			outString.WriteString(fmt.Sprintf("%v> %v\n", strings.Repeat(" ", spaces), redirect.Headers.Get("Location")))
			out.WriteString(outString.String())
		}
	} else {
		fmt.Println("No redirects")
	}
	fmt.Print(out.String())

}

func serverFingerprint(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	if !fc.ShowFp {
		return
	}
	out.WriteString(types.Cyan + "\nServer Fingerprint:\n" + types.Reset)
	out.WriteString("   Server: " + bl.Headers.Get("Server") + "\n")
	if bl.Headers.Get("X-Powered-By") != "" {
		out.WriteString("   Tech: " + bl.Headers.Get("X-Powered-By") + "\n")

	}
	if bl.Headers.Get("X-Powered-CMS") != "" {
		out.WriteString(types.Cyan + "   CMS: " + types.Reset + bl.Headers.Get("X-Powered-CMS") + "\n")

	}
	if bl.Headers.Get("X-Frame-Options") == "" {
		out.WriteString(types.Cyan + "   [INFO] " + types.Reset + "Missing X-Frame-Options header. (Clickjacking-related behavior)\n")
	}
	out.WriteString(types.Cyan + "Defined links: \n" + types.Reset)
	out.WriteString(getLinks(bl))

	fmt.Println(out.String())

}

func getLinks(bl types.BlinkResponse) string {
	var out strings.Builder
	match, _ := regexp.MatchString(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`, string(bl.Body))
	if !match {
		return "No links in body."
	}
	r, _ := regexp.Compile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`)
	links := r.FindAllString(string(bl.Body), -1)
	for _, link := range links {
		out.WriteString(fmt.Sprintf("   %s\n", link))
	}
	return out.String()
}

func bodyOutput(bl types.BlinkResponse, fc types.FlagCondition) string {
	if fc.ShowBody {
		return bl.BodyPreview
	} else if fc.ShowFullBody {
		return string(bl.Body)
	}
	return ""
}

func ErrorOutput(err types.BlinkError) string {
	var out strings.Builder
	if err.Stage != "OK" {
		if err.Stage == "Unknown" {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ] %s\n"+types.Reset, err.Stage, err.Message))
		} else if err.Stage == "INFO" {
			out.WriteString(fmt.Sprintf(types.Yellow+"[ %v ] %v"+types.Reset, err.Stage, err.Message))
		} else {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ] %v \n"+types.Reset, err.Stage, err.Message))
		}

	}
	return out.String()
}

func Diffs(bl []types.BlinkResponse, fc types.FlagCondition) {
	var baseline = bl[0]

	if baseline.URL == "" {
		fmt.Println("EMPTY BASELINE")
		return
	}

	for _, r := range bl[1:] {
		var out strings.Builder
		hasDiff := false
		if r.URL == "" {
			fmt.Println("EMPTY REQUEST")
			continue
		}

		if baseline.StatusCode != r.StatusCode {
			diffLine(&out, "status",
				fmt.Sprintf("%d -> %d", baseline.StatusCode, r.StatusCode), fc,
			)
			hasDiff = true
		}

		if baseline.BodyHash != r.BodyHash {
			diffLine(&out, "body_hash",
				fmt.Sprintf("%v.. -> %v..", shortHash(baseline.BodyHash), shortHash(r.BodyHash)), fc,
			)
			hasDiff = true

			parts := strings.FieldsFunc(r.RequestData, func(g rune) bool {
				return g == '=' || g == '&'
			})

			if len(parts)%2 == 0 {
				for i := 1; i < len(parts); i += 2 {
					value := parts[i]
					if strings.Contains(string(r.Body), value) {
						diffLine(&out, "reflect", "raw input", fc)
						hasDiff = true

					}
					if strings.Contains(string(r.Body), url.QueryEscape(value)) || strings.Contains(string(r.Body), url.QueryEscape(url.QueryEscape(value))) {

						diffLine(&out, "reflect", "encoded input", fc)
						hasDiff = true

					}
				}
			}
		}
		headersChanges := diffHeaders(baseline.Headers, r.Headers)
		if len(headersChanges) > 0 {
			diffLine(&out, "headers", strings.Join(headersChanges, ", "), fc)
			hasDiff = true
		}

		if baseline.Timings.FullRtt*2 < r.Timings.FullRtt {
			diffLine(&out, "rtt",
				fmt.Sprintf("%v -> %v",
					baseline.Timings.FullRtt,
					r.Timings.FullRtt,
				), fc,
			)
			hasDiff = true
		}

		if len(baseline.Cookies) != len(r.Cookies) {
			diffLine(&out, "cookies",
				fmt.Sprintf("%d -> %d",
					len(baseline.Cookies),
					len(r.Cookies),
				), fc,
			)
			hasDiff = true
		}
		if hasDiff {
			if r.Method == "POST" {
				fmt.Printf("  %s", types.Cyan+r.RawRequest.URL.Path+types.Reset+" "+r.RequestData+"\n")
			} else {
				fmt.Printf("  %s", types.Cyan+r.RawRequest.URL.Path+types.Reset+" "+r.RawRequest.URL.RawQuery+types.Reset+" "+r.RequestData+"\n")
			}

		}
		fmt.Println(out.String())
	}

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
func diffLine(out *strings.Builder, field string, msg string, fc types.FlagCondition) {
	if fc.IgnoreHash {
		if field == "body_hash" {
			return
		}
	}
	if fc.IgnoreReflection {
		if field == "reflect" {
			return
		}
	}
	out.WriteString(fmt.Sprintf(
		types.Magenta+"    %s"+types.Reset+" : %s\n",
		field,
		msg,
	))
}

func shortHash(hash string) string {
	if hash != "" {
		return hash[:10]
	} else {
		return "EMPTY_HASH"
	}
}
