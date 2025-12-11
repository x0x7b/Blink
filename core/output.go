package core

import (
	"fmt"
	"strings"
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

func CleanOutput(bl BlinkResponse, mode int, fc FlagCondition) {
	var out strings.Builder
	if mode == 0 {

		out.WriteString(Cyan + "[Blink] ===================================================================================\n" + Reset)
		out.WriteString(fmt.Sprintf(
			Bold+"method: "+Reset+"%s\n"+
				Bold+"url:    "+Reset+"%s\n"+
				Bold+"status: "+Reset+ColorStatus(bl.StatusCode)+"%d (%s)"+Reset+"\n"+
				Bold+"proto:  "+Reset+"%s (%d.%d)\n"+
				Bold+"rtt:    "+Reset+"%s\n"+
				Bold+"length: "+Reset+"%d\n",
			bl.Method,
			bl.URL,
			bl.StatusCode, bl.Status,
			bl.Proto, bl.ProtoMajor, bl.ProtoMinor,
			bl.RTT,
			bl.ContentLength,
		))
		out.WriteString(Bold + "TLS:" + Reset + "\n")
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

	} else {
		out.WriteString(ColorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + Reset)
		out.WriteString(Blue + "[ " + Reset + Cyan + bl.Method + Reset + " " + bl.URL + Blue + " ] " + Reset)
		out.WriteString(fmt.Sprintf("(%vms)\n", bl.RTT.Milliseconds()))

	}
	fmt.Print(out.String())

}
