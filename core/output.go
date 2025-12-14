package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

func ColorStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return Green
	case code >= 300 && code < 400:
		return Blue
	case code >= 400 && code < 500:
		return Yellow
	default:
		return Red
	}
}

func colorTime(timing time.Duration) string {
	switch {
	case timing < 150*time.Millisecond:
		return Green
	case timing >= 150*time.Microsecond:
		return Yellow
	case timing >= 300*time.Millisecond:
		return Red
	default:
		return Red
	}

}

func CleanOutput(bl BlinkResponse, rc []BlinkResponse, fc FlagCondition) {
	fmt.Println(Magenta + "[ Blink v0.4 ]  \n" + Reset)
	if len(rc) > 0 {
		redirectChainOutput(rc, fc)
	}
	switch {
	case fc.OutputMode == 0:
		fmt.Printf(Cyan + "Final response:\n" + Reset)
		defaultOutput(bl, fc)
	case fc.OutputMode == 1:
		verboseOutput(bl, fc)
	}

}

func defaultOutput(bl BlinkResponse, fc FlagCondition) {
	var out strings.Builder
	out.WriteString(ColorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + Reset)
	out.WriteString(Blue + "[ " + Reset + Cyan + bl.Method + Reset + " " + bl.URL + Blue + " ] " + Reset)
	out.WriteString(fmt.Sprintf("(%v)\n", bl.Timings.fullRtt))
	if fc.ShowBody {
		out.WriteString("\n")
		out.WriteString(bl.BodyPreview)
		out.WriteString("\n")
	}
	fmt.Print(out.String())
	out.WriteString(bodyOutput(bl, fc))
	serverFingerprint(bl)

}

func verboseOutput(bl BlinkResponse, fc FlagCondition) {
	var out strings.Builder
	out.WriteString(Cyan + Bold + "Final response: \n" + Reset)
	out.WriteString(fmt.Sprintf(
		Bold+"method: "+Reset+"%s\n"+
			Bold+"url:    "+Reset+"%s\n"+
			Bold+"status: "+Reset+ColorStatus(bl.StatusCode)+"%d (%s)"+Reset+"\n"+
			Bold+"proto:  "+Reset+"%s (%d.%d)\n"+
			Bold+"rtt:    "+Reset+"%s\n"+
			Bold+"  dns:    "+Reset+colorTime(bl.Timings.dnsDuration)+"%s\n"+Reset+
			Bold+"  tcp:    "+Reset+colorTime(bl.Timings.tcpDuration)+"%s\n"+Reset+
			Bold+"  tls:    "+Reset+colorTime(bl.Timings.tlsDuration)+"%s\n"+Reset+
			Bold+"  ttfb:   "+Reset+colorTime(bl.Timings.ttfb)+"%s\n"+Reset+
			Bold+"  length: "+Reset+"%d\n",
		bl.Method,
		bl.URL,
		bl.StatusCode, bl.Status,
		bl.Proto, bl.ProtoMajor, bl.ProtoMinor,
		bl.Timings.fullRtt,
		bl.Timings.dnsDuration,
		bl.Timings.tcpDuration,
		bl.Timings.tlsDuration,
		bl.Timings.ttfb,
		bl.ContentLength,
	))

	out.WriteString(Cyan + "TLS:" + Reset + "\n")
	out.WriteString(fmt.Sprintf("   alpn: %v\n", bl.ALPN))
	out.WriteString(fmt.Sprintf("   Version: %v\n", bl.TLSVersion))
	out.WriteString(fmt.Sprintf("   Cipher:  %v\n", bl.CipherSuite))
	out.WriteString(fmt.Sprintf("   Issuer:  %v\n", bl.CertIssuer))
	out.WriteString(fmt.Sprintf("   Expires: %v\n", bl.CertExpires))

	out.WriteString(Cyan + "headers:" + Reset + "\n")
	for k, v := range bl.Headers {
		out.WriteString(fmt.Sprintf("  "+Cyan+"%s:"+Reset+" %s\n", k, v))
	}

	out.WriteString(bodyOutput(bl, fc))
	fmt.Print(out.String())
	serverFingerprint(bl)
}

func redirectChainOutput(redirects []BlinkResponse, fc FlagCondition) {
	var out strings.Builder
	var stringWidth int
	out.WriteString(Cyan + "Redirect chain:\n" + Reset)
	if len(redirects) > 0 {
		var maxLenRequest int
		for i, req := range redirects {
			var outForLen strings.Builder
			outForLen.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outForLen.WriteString(fmt.Sprintf("%v ", req.StatusCode))
			outForLen.WriteString("[ " + req.Method + " " + req.URL + " ] ")
			outForLen.WriteString(fmt.Sprintf("(%v)", req.Timings.fullRtt))
			if len(outForLen.String()) > maxLenRequest {
				maxLenRequest = len(outForLen.String())
			}
			stringWidth = maxLenRequest + 1
		}

		for i, redirect := range redirects {
			var outString strings.Builder
			var outNoColors strings.Builder
			outString.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outString.WriteString(ColorStatus(redirect.StatusCode) + fmt.Sprintf("%v ", redirect.StatusCode) + Reset)
			outString.WriteString(Blue + "[ " + Reset + Cyan + redirect.Method + Reset + " " + redirect.URL + Blue + " ] " + Reset)
			outString.WriteString(fmt.Sprintf("(%v)", redirect.Timings.fullRtt))

			outNoColors.WriteString(fmt.Sprintf("   [ %d ] ", i))
			outNoColors.WriteString(fmt.Sprintf("%v ", redirect.StatusCode))
			outNoColors.WriteString("[ " + redirect.Method + " " + redirect.URL + " ] ")
			outNoColors.WriteString(fmt.Sprintf("(%v)", redirect.Timings.fullRtt))

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

func serverFingerprint(bl BlinkResponse) {
	var out strings.Builder
	out.WriteString(Cyan + "\nServer Fingerprint:\n" + Reset)
	out.WriteString("   Server: " + bl.Headers.Get("Server") + "\n")
	if bl.Headers.Get("X-Powered-By") != "" {
		out.WriteString("   Tech: " + bl.Headers.Get("X-Powered-By") + "\n")

	}
	if bl.Headers.Get("X-Powered-CMS") != "" {
		out.WriteString(Cyan + "   CMS: " + Reset + bl.Headers.Get("X-Powered-CMS") + "\n")

	}
	if bl.Headers.Get("X-Frame-Options") == "" {
		out.WriteString("   Missing X-Frame-Options header. (Clickjacking-related behavior)")
	}
	out.WriteString(Cyan + "Defined links: \n" + Reset)
	out.WriteString(getLinks(bl))
	fmt.Println(out.String())

}

func getLinks(bl BlinkResponse) string {
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

func bodyOutput(bl BlinkResponse, fc FlagCondition) string {
	if fc.ShowBody {
		return bl.BodyPreview
	} else if fc.ShowFullBody {
		return string(bl.Body)
	}
	return ""
}
