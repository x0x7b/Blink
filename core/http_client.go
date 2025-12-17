package core

import (
	"Blink/types"

	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

func HttpRequest(method string, domain string, fc types.FlagCondition) (types.BlinkResponse, []types.BlinkResponse, types.BlinkError) {
	current := domain
	currentMethod := method
	var redirectChain []types.BlinkResponse

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // НЕ редіректитись
		},
		Timeout: time.Duration(fc.Timeout) * time.Second,
	}

	for hop := 0; hop < fc.MaxRedirects; hop++ {
		var networkTimings types.NetworkTimings
		var blinkResp types.BlinkResponse
		var dnsStart time.Time
		var dnsDuration time.Duration
		var connectStart time.Time
		var tcpDuration time.Duration
		var tlsStart time.Time
		var tlsDuration time.Duration
		var start time.Time
		var ttfb time.Duration

		trace := &httptrace.ClientTrace{
			DNSStart: func(info httptrace.DNSStartInfo) {
				dnsStart = time.Now()
			},
			DNSDone: func(info httptrace.DNSDoneInfo) {
				if !dnsStart.IsZero() {
					dnsDuration = time.Since(dnsStart)
				}
			},

			ConnectStart: func(network, addr string) {
				connectStart = time.Now()
			},
			ConnectDone: func(network, addr string, err error) {
				if err == nil {
					tcpDuration = time.Since(connectStart)
				}
			},

			TLSHandshakeStart: func() {
				tlsStart = time.Now()
			},
			TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
				tlsDuration = time.Since(tlsStart)
			},

			GotFirstResponseByte: func() {
				ttfb = time.Since(start)
			},
		}

		if currentMethod == "GET" && fc.Data != "" {
			return blinkResp, redirectChain, ClassifyNetworkError(fmt.Errorf("GET request cannot have a body; use -X POST or -X PUT"))
		}

		var payloadReader io.Reader
		if fc.Data != "" {
			payloadReader = strings.NewReader(fc.Data)
		} else {
			payloadReader = nil
		}

		req, err := http.NewRequest(currentMethod, current, payloadReader)
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
		if err != nil {
			return blinkResp, redirectChain, ClassifyNetworkError(err)
		}

		if fc.Data != "" {
			req.Header.Set("Content-Type", "application/json")

		}

		req.Header.Set("User-Agent", "Blink/1.0")
		req.Header.Set("Accept", "*/*")
		if req.Method == "POST" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		}
		start = time.Now()

		resp, err := client.Do(req)

		networkTimings.FullRtt = time.Since(start)
		networkTimings.DnsDuration = dnsDuration
		networkTimings.TcpDuration = tcpDuration
		networkTimings.TlsDuration = tlsDuration
		networkTimings.Ttfb = ttfb

		if err != nil {
			return blinkResp, redirectChain, ClassifyNetworkError(err)
		}

		blinkResp, err = makeBlinkResponce(resp, networkTimings)
		if err != nil {
			return blinkResp, redirectChain, ClassifyNetworkError(err)
		}
		blinkResp.RequestData = fc.Data
		resp.Body.Close()
		if !fc.FollowRedirects {
			return blinkResp, redirectChain, ClassifyNetworkError(fmt.Errorf("redirect received but not followed (--no-follow)"))
		}

		if blinkResp.StatusCode >= 300 && blinkResp.StatusCode < 400 { // redirect
			loc := blinkResp.Headers.Get("Location")
			if loc == "" {
				return blinkResp, redirectChain, ClassifyNetworkError(fmt.Errorf("redirect with no Location header"))
			}

			u, err := url.Parse(loc)
			if err != nil {
				return blinkResp, redirectChain, ClassifyNetworkError(fmt.Errorf("invalid redirect Location"))
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
			return blinkResp, redirectChain, ClassifyNetworkError(nil)
		}

	}
	return types.BlinkResponse{}, redirectChain, ClassifyNetworkError(nil)
}

func makeBlinkResponce(resp *http.Response, timings types.NetworkTimings) (types.BlinkResponse, error) {
	var blinkResp types.BlinkResponse

	blinkResp.RawRequest = resp.Request
	blinkResp.RawResponse = resp

	blinkResp.Status = resp.Status
	blinkResp.StatusCode = resp.StatusCode
	blinkResp.Proto = resp.Proto
	blinkResp.ProtoMajor = resp.ProtoMajor
	blinkResp.ProtoMinor = resp.ProtoMinor
	blinkResp.Headers = resp.Header
	blinkResp.ContentLength = resp.ContentLength
	blinkResp.Method = resp.Request.Method
	blinkResp.URL = resp.Request.URL.String()
	blinkResp.Timings = timings
	blinkResp.Cookies = resp.Cookies()
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		blinkResp.Redirected = true
	}

	// TLS

	if resp.TLS != nil {
		blinkResp.ALPN = resp.TLS.NegotiatedProtocol
	}

	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		blinkResp.TLSVersion = resp.TLS.Version
		blinkResp.CipherSuite = resp.TLS.CipherSuite
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
