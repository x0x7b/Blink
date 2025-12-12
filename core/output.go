package core

import (
	"fmt"
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

func CleanOutput(bl BlinkResponse, rc []BlinkResponse, mode int, fc FlagCondition) {
	switch {
	case mode == 0:
		defaultOutput(bl, fc)
	case mode == 1:
		verboseOutput(bl, fc)
	case mode == 2:
		redirectChainOutput(rc, fc)
	}
}

func defaultOutput(bl BlinkResponse, fc FlagCondition) {
	var out strings.Builder
	out.WriteString(ColorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + Reset)
	out.WriteString(Blue + "[ " + Reset + Cyan + bl.Method + Reset + " " + bl.URL + Blue + " ] " + Reset)
	out.WriteString(fmt.Sprintf("(%vms)\n", bl.Timings.fullRtt))
	fmt.Print(out.String())

}

func verboseOutput(bl BlinkResponse, fc FlagCondition) {
	var out strings.Builder
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

	out.WriteString(Bold + "TLS:" + Reset + "\n")
	out.WriteString(fmt.Sprintf("   alpn: %v\n", bl.ALPN))
	out.WriteString(fmt.Sprintf("   Version: %v\n", bl.TLSVersion))
	out.WriteString(fmt.Sprintf("   Cipher:  %v\n", bl.CipherSuite))
	out.WriteString(fmt.Sprintf("   Issuer:  %v\n", bl.CertIssuer))
	out.WriteString(fmt.Sprintf("   Expires: %v\n", bl.CertExpires))

	out.WriteString(Bold + "headers:" + Reset + "\n")
	for k, v := range bl.Headers {
		out.WriteString(fmt.Sprintf("  "+Bold+"%s:"+Reset+" %s\n", k, v))
	}

	if fc.ShowBody {
		out.WriteString("\n")
		out.WriteString(bl.BodyPreview)
		out.WriteString("\n")
	}

	fmt.Print(out.String())
}

func redirectChainOutput(redirects []BlinkResponse, fc FlagCondition) {
	var out strings.Builder
	if len(redirects) > 0 {
		for i, redirect := range redirects {
			out.WriteString("[ " + string(i) + " ]")
			out.WriteString(ColorStatus(redirect.StatusCode) + fmt.Sprintf("%v ", redirect.StatusCode) + Reset)
			out.WriteString(Blue + "[ " + Reset + Cyan + redirect.Method + Reset + " " + redirect.URL + Blue + " ] " + Reset)
			out.WriteString(fmt.Sprintf("(%vms)\n", redirect.Timings.fullRtt))
		}
	}

}
