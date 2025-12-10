package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BlinkResponse struct {
	StatusCode    int
	Status        string
	Proto         string
	ProtoMajor    int
	ProtoMinor    int
	Headers       http.Header
	Body          []byte
	BodyPreview   string
	BodyHash      string
	ContentLength int64
	RTT           time.Duration
	Method        string
	URL           string
	TLSVersion    uint16
	CipherSuite   uint16
	CertIssuer    string
	CertExpires   time.Time
}

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Blue   = "\033[34m"
	White  = "\033[37m"
)

func colorStatus(code int) string {
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

type flagCondition struct {
	ShowBody        bool
	FollowRedirects bool
	MaxRedirects    int
}

var fc flagCondition

func main() {

	showBody := flag.Bool("b", false, "Show response body")
	showBodyLong := flag.Bool("include-body", false, "Show response body")

	method := flag.String("X", "GET", "HTTP method")

	followRedirects := flag.Bool("no-follow", true, "Follow redirects")
	maxRedirects := flag.Int("max-redirects", 5, "Set value to max redirects")

	data := flag.String("data", "", "HTTP Payload")
	flag.Parse()
	if *showBody || *showBodyLong {
		fc.ShowBody = true
	}
	fc.FollowRedirects = *followRedirects
	fc.MaxRedirects = *maxRedirects

	if flag.NArg() < 1 {
		log.Fatal("URL is required")
	}

	url := flag.Arg(0)
	response, redirects, err := HttpRequest(*method, url, *data)
	if err != nil {
		log.Println(err)
	}
	if len(redirects) > 0 {
		for _, redirect := range redirects {
			cleanOutput(redirect, 0)
		}
	}
	cleanOutput(response, 1)
}

func HttpRequest(method string, domain string, data string) (BlinkResponse, []BlinkResponse, error) {
	current := domain
	currentMethod := method
	var redirectChain []BlinkResponse

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // НЕ редіректитись
		},
	}

	for hop := 0; hop < fc.MaxRedirects; hop++ {
		var blinkResp BlinkResponse

		if currentMethod == "GET" && data != "" {
			return blinkResp, redirectChain, fmt.Errorf("GET request cannot have a body; use -X POST or -X PUT")
		}

		var payloadReader io.Reader
		if data != "" {
			payloadReader = strings.NewReader(data)
		} else {
			payloadReader = nil
		}

		req, err := http.NewRequest(currentMethod, current, payloadReader)

		if err != nil {
			return blinkResp, redirectChain, err
		}

		if data != "" {
			req.Header.Set("Content-Type", "application/json")

		}

		req.Header.Set("User-Agent", "Blink/1.0")
		req.Header.Set("Accept", "*/*")
		start := time.Now()
		resp, err := client.Do(req)
		rtt := time.Since(start)

		if err != nil {
			return blinkResp, redirectChain, err
		}

		blinkResp, err = makeBlinkResponce(resp, rtt)
		if err != nil {
			return blinkResp, redirectChain, err
		}
		resp.Body.Close()
		if !fc.FollowRedirects {
			return blinkResp, redirectChain, fmt.Errorf("handeled redirect, but deny from user flags")
		}

		if blinkResp.StatusCode >= 300 && blinkResp.StatusCode < 400 { // redirect
			loc := blinkResp.Headers.Get("Location")
			if loc == "" {
				return blinkResp, redirectChain, fmt.Errorf("redirect with no Location header")
			}

			u, err := url.Parse(loc)
			if err != nil {
				return blinkResp, redirectChain, fmt.Errorf("invalid redirect Location")
			}
			if blinkResp.StatusCode == 302 || blinkResp.StatusCode == 303 {
				currentMethod = "GET"
				payloadReader = nil
			} else {

			}
			current = resp.Request.URL.ResolveReference(u).String()
			redirectChain = append(redirectChain, blinkResp)
			continue

		}
		if blinkResp.StatusCode < 300 || blinkResp.StatusCode >= 400 { // not redirect
			return blinkResp, redirectChain, nil
		}

	}
	return BlinkResponse{}, redirectChain, nil
}

func cleanOutput(bl BlinkResponse, mode int) {
	var out strings.Builder
	if mode == 0 {

		out.WriteString(Cyan + "[Blink] ===================================================================================\n" + Reset)
		out.WriteString(fmt.Sprintf(
			Bold+"method: "+Reset+"%s\n"+
				Bold+"url:    "+Reset+"%s\n"+
				Bold+"status: "+Reset+colorStatus(bl.StatusCode)+"%d (%s)"+Reset+"\n"+
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
		out.WriteString(colorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + Reset)
		out.WriteString(Blue + "[ " + Reset + bl.URL + Blue + " ] " + Reset)
		out.WriteString(fmt.Sprintf("(%vms)\n", bl.RTT.Milliseconds()))

	}
	fmt.Print(out.String())

}

func makeBlinkResponce(resp *http.Response, rtt time.Duration) (BlinkResponse, error) {
	var blinkResp BlinkResponse
	blinkResp.Status = resp.Status
	blinkResp.StatusCode = resp.StatusCode
	blinkResp.Proto = resp.Proto
	blinkResp.ProtoMajor = resp.ProtoMajor
	blinkResp.ProtoMinor = resp.ProtoMinor
	blinkResp.Headers = resp.Header
	blinkResp.ContentLength = resp.ContentLength
	blinkResp.Method = resp.Request.Method
	blinkResp.URL = resp.Request.URL.String()
	blinkResp.RTT = rtt

	// TLS
	blinkResp.TLSVersion = resp.TLS.Version
	blinkResp.CipherSuite = resp.TLS.CipherSuite
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0] // leaf certificate
		blinkResp.CertIssuer = cert.Issuer.String()
		blinkResp.CertExpires = cert.NotAfter
	}

	limited := io.LimitReader(resp.Body, 2*1024*1024) // 2MB
	bodyBytes, _ := io.ReadAll(limited)
	blinkResp.Body = bodyBytes
	if len(bodyBytes) > 300 {
		blinkResp.BodyPreview = string(bodyBytes[:300])
	} else {
		blinkResp.BodyPreview = string(bodyBytes)
	}

	sum := sha256.Sum256(bodyBytes)
	blinkResp.BodyHash = fmt.Sprintf("%x", sum)

	return blinkResp, nil
}
