package output

import (
	"Blink/core"
	"Blink/types"
	"fmt"
	"strings"
)

func CleanOutput(bl types.BlinkResponse, rc []types.BlinkResponse, fc types.FlagCondition) {
	if len(rc) > 0 && bl.URL != "" {
		redirectChainOutput(rc, fc)
	}
	if len(rc) > 0 && bl.URL == "" {
		core.Diffs(rc, fc)
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
	out.WriteString(fmt.Sprintf("(%v%v%v)\n", colorTime(bl.Timings.FullRtt), bl.Timings.FullRtt, types.Reset))
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
		out.WriteString(fmt.Sprintf("  %-12s:"+types.Reset+" %s\n", k, v))
	}

	out.WriteString(bodyOutput(bl, fc))
	fmt.Print(out.String())
	serverFingerprint(bl, fc)
}
