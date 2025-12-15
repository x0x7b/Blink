package core

import (
	"Blink/types"
	"fmt"
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
	if len(rc) > 0 {
		redirectChainOutput(rc, fc)
	}
	switch {
	case fc.OutputMode == 0:
		fmt.Printf(types.Cyan + "Final response:\n" + types.Reset)
		defaultOutput(bl, fc)
	case fc.OutputMode == 1:
		verboseOutput(bl, fc)
	}

}

func defaultOutput(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	out.WriteString(ColorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + types.Reset)
	out.WriteString(types.Blue + "[ " + types.Reset + types.Cyan + bl.Method + types.Reset + " " + bl.URL + types.Blue + " ] " + types.Reset)
	out.WriteString(fmt.Sprintf("(%v)\n", bl.Timings.FullRtt))
	if fc.ShowBody {
		out.WriteString("\n")
		out.WriteString(bl.BodyPreview)
		out.WriteString("\n")
	}
	fmt.Print(out.String())
	out.WriteString(bodyOutput(bl, fc))
	serverFingerprint(bl, fc)

}

func verboseOutput(bl types.BlinkResponse, fc types.FlagCondition) {
	var out strings.Builder
	out.WriteString(types.Cyan + types.Bold + "Final response: \n" + types.Reset)
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

func errorOutput(err types.BlinkError) string {
	var out strings.Builder
	if err.Stage != "OK" {
		if err.Stage == "Unknown" {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ] %s\n"+types.Reset, err.Stage, err.Message))
		} else if err.Stage == "INFO" {
			out.WriteString(fmt.Sprintf(types.Yellow+"[ %v ] %v"+types.Reset, err.Stage, err.Message))
		} else {
			out.WriteString(fmt.Sprintf(types.Red+"[ %v ERROR ]\n"+types.Reset, err.Stage))
		}

	}
	return out.String()
}
