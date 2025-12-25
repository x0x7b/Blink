package types

import (
	"net/http"
	"time"
)

type FlagCondition struct {
	ShowBody         bool
	FollowRedirects  bool
	MaxRedirects     int
	Timeout          int
	OutputMode       int
	ShowFullBody     bool
	Data             string
	TestParam        bool
	TestForms        bool
	ShowFp           bool
	IgnoreHash       bool
	IgnoreReflection bool
	Wordlist         string
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
	Cookies       []*http.Cookie
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
	Redirected    bool
	RequestData   string
}

type BlinkError struct {
	Stage   string // DNS, TCP, TLS, HTTP, REDIRECT, BODY, UNKNOWN, INFO, OK
	Message string
}

type NetworkTimings struct {
	DnsDuration time.Duration
	TcpDuration time.Duration
	TlsDuration time.Duration
	Ttfb        time.Duration
	FullRtt     time.Duration
}

type Form struct {
	Name   string
	Method string
	Action string
	Inputs []Input
}

type Input struct {
	Name string
	Type string
}

type Progress struct {
	Stage   string
	Current int
	Total   int
	Target  string
}

type BehaviorProfile struct {
	TotalTests int
	Ratios     map[DiffKind]float64
	Counts     map[DiffKind]int
}

type Diff struct {
	Kind   DiffKind
	Before string
	After  string
	Meta   string
}

type TestResult struct {
	Payload string
	Score   float64
	Path    string
	Diffs   []Diff
}

type DiffKind int

const (
	DiffBodyHash DiffKind = iota
	DiffStatus
	DiffRTT
	DiffReflect
	DiffHeaders
	DiffCookies
)

func (k DiffKind) String() string {
	switch k {
	case DiffBodyHash:
		return "BODY_HASH"
	case DiffStatus:
		return "STATUS"
	case DiffRTT:
		return "RTT"
	case DiffReflect:
		return "REFLECT"
	case DiffHeaders:
		return "HEADERS"
	case DiffCookies:
		return "COOKIES"
	default:
		return "UNKNOWN"
	}
}

func (k DiffKind) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

var BaseWeight = map[DiffKind]float64{
	DiffBodyHash: 1.0,
	DiffStatus:   1.0,
	DiffRTT:      1.0,
	DiffReflect:  1.0,
	DiffHeaders:  1.0,
	DiffCookies:  1.0,
}
