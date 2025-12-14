package core

import (
	"net/http"
	"time"
)

type FlagCondition struct {
	ShowBody        bool
	FollowRedirects bool
	MaxRedirects    int
	Timeout         int
	OutputMode      int
	ShowFullBody    bool
	Data            string
}

const (
	Reset   = "\033[0m"
	Bold    = "\033[1m"
	Cyan    = "\033[36m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Red     = "\033[31m"
	Blue    = "\033[34m"
	White   = "\033[37m"
	Magenta = "\033[35m"
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
	Timings       NetworkTimings
	ALPN          string
	RawRequest    *http.Request
	RawResponse   *http.Response
}

type BlinkError struct {
	Stage   string // DNS, TCP, TLS, HTTP, REDIRECT, BODY, UNKNOWN, INFO, OK
	Message string
}

type NetworkTimings struct {
	dnsDuration time.Duration
	tcpDuration time.Duration
	tlsDuration time.Duration
	ttfb        time.Duration
	fullRtt     time.Duration
}
