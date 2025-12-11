package core

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpRequest(method string, domain string, data string, fc FlagCondition) (BlinkResponse, []BlinkResponse, BlinkError) {
	current := domain
	currentMethod := method
	var redirectChain []BlinkResponse

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // НЕ редіректитись
		},
		Timeout: time.Duration(fc.Timeout) * time.Second,
	}

	for hop := 0; hop < fc.MaxRedirects; hop++ {
		var blinkResp BlinkResponse

		if currentMethod == "GET" && data != "" {
			return blinkResp, redirectChain, classifyNetworkError(fmt.Errorf("GET request cannot have a body; use -X POST or -X PUT"))
		}

		var payloadReader io.Reader
		if data != "" {
			payloadReader = strings.NewReader(data)
		} else {
			payloadReader = nil
		}

		req, err := http.NewRequest(currentMethod, current, payloadReader)

		if err != nil {
			return blinkResp, redirectChain, classifyNetworkError(err)
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
			return blinkResp, redirectChain, classifyNetworkError(err)
		}

		blinkResp, err = makeBlinkResponce(resp, rtt)
		if err != nil {
			return blinkResp, redirectChain, classifyNetworkError(err)
		}
		resp.Body.Close()
		if !fc.FollowRedirects {
			return blinkResp, redirectChain, classifyNetworkError(fmt.Errorf("handeled redirect, but deny from user flags"))
		}

		if blinkResp.StatusCode >= 300 && blinkResp.StatusCode < 400 { // redirect
			loc := blinkResp.Headers.Get("Location")
			if loc == "" {
				return blinkResp, redirectChain, classifyNetworkError(fmt.Errorf("redirect with no Location header"))
			}

			u, err := url.Parse(loc)
			if err != nil {
				return blinkResp, redirectChain, classifyNetworkError(fmt.Errorf("invalid redirect Location"))
			}
			if blinkResp.StatusCode == 302 || blinkResp.StatusCode == 303 || blinkResp.StatusCode == 301 {
				currentMethod = "GET"
				payloadReader = nil
			} else {

			}
			current = resp.Request.URL.ResolveReference(u).String()
			redirectChain = append(redirectChain, blinkResp)
			continue

		}
		if blinkResp.StatusCode < 300 || blinkResp.StatusCode >= 400 { // not redirect
			return blinkResp, redirectChain, classifyNetworkError(nil)
		}

	}
	return BlinkResponse{}, redirectChain, classifyNetworkError(nil)
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
