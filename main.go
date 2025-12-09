package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	Reset = "\033[0m"
	Bold  = "\033[1m"

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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ")
	}
	responce, err := HttpRequest(os.Args[1])
	if err != nil {
		log.Println(err)
	}
	cleaOutput(responce)
}

func HttpRequest(domain string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", domain, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Blink/1.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return resp, nil
}

func cleaOutput(responce *http.Response) {
	status := responce.Status
	statusCode := responce.StatusCode
	proto := responce.Proto
	protoMajor := responce.ProtoMajor
	protoMinor := responce.ProtoMinor
	header := responce.Header
	contentLength := responce.ContentLength
	method := responce.Request.Method
	url := responce.Request.URL

	fmt.Printf(
		Cyan+"[Blink]"+Reset+"\n"+
			Bold+"method:         "+Reset+"%s\n"+
			Bold+"url:            "+Reset+"%s\n"+
			Bold+"status:         "+Reset+colorStatus(statusCode)+"%d (%s)"+Reset+"\n"+
			Bold+"proto:          "+Reset+"%s (%d.%d)\n"+
			Bold+"content_length: "+Reset+"%d\n"+
			Bold+"headers:"+Reset+"\n",
		method,
		url,
		statusCode, status,
		proto, protoMajor, protoMinor,
		contentLength,
	)

	for k, v := range header {
		fmt.Printf("  "+Bold+"%s:"+Reset+" %s\n", k, v)
	}

}
