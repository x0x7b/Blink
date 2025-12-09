package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
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
	ShowBody bool
}

var fc flagCondition

func main() {

	showBody := flag.Bool("b", false, "Show response body")
	showBodyLong := flag.Bool("include-body", false, "Show response body")

	method := flag.String("X", "GET", "HTTP method")

	data := flag.String("data", "", "HTTP Payload")
	flag.Parse()
	if *showBody || *showBodyLong {
		fc.ShowBody = true
	}

	if flag.NArg() < 1 {
		log.Fatal("URL is required")
	}

	url := flag.Arg(0)
	response, err := HttpRequest(*method, url, *data)
	if err != nil {
		log.Println(err)
	}
	cleanOutput(response, 1)
}

func HttpRequest(method string, domain string, data string) (BlinkResponse, error) {
	var blinkResp BlinkResponse
	if method == "GET" && data != "" {
		return blinkResp, fmt.Errorf("GET request cannot have a body; use -X POST or -X PUT")
	}
	var payloadReader io.Reader
	if data != "" {
		payloadReader = strings.NewReader(data)
	} else {
		payloadReader = nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, domain, payloadReader)
	if err != nil {
		return blinkResp, err
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
		return blinkResp, err
	}
	defer resp.Body.Close()

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

	limited := io.LimitReader(resp.Body, 2*1024*1024) // 2MB
	bodyBytes, _ := io.ReadAll(limited)
	blinkResp.Body = bodyBytes
	if len(bodyBytes) > 300 {
		blinkResp.BodyPreview = string(bodyBytes[:300])
	} else {
		blinkResp.BodyPreview = string(bodyBytes)
	}

	return blinkResp, nil
}

func cleanOutput(bl BlinkResponse, mode int) {
	if mode == 0 {
		var out strings.Builder

		out.WriteString(Cyan + "[Blink]\n" + Reset)
		out.WriteString(fmt.Sprintf(
			Bold+"method:         "+Reset+"%s\n"+
				Bold+"url:            "+Reset+"%s\n"+
				Bold+"status:         "+Reset+colorStatus(bl.StatusCode)+"%d (%s)"+Reset+"\n"+
				Bold+"proto:          "+Reset+"%s (%d.%d)\n"+
				Bold+"rtt:            "+Reset+"%s\n"+
				Bold+"content_length: "+Reset+"%d\n",
			bl.Method,
			bl.URL,
			bl.StatusCode, bl.Status,
			bl.Proto, bl.ProtoMajor, bl.ProtoMinor,
			bl.RTT,
			bl.ContentLength,
		))
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
	} else {
		var out strings.Builder
		out.WriteString(colorStatus(bl.StatusCode) + fmt.Sprintf("%v ", bl.StatusCode) + Reset)
		out.WriteString(Blue + "[ " + Reset + bl.URL + Blue + " ] " + Reset)
		out.WriteString(fmt.Sprintf("(%vms)", bl.RTT.Milliseconds()))

		fmt.Print(out.String())
	}

}
