package core

import (
	"Blink/types"
	"crypto/x509"
	"errors"
	"net"
)

func ClassifyNetworkError(err error) types.BlinkError {
	var be types.BlinkError

	// nil error
	if err == nil {
		be.Stage = "OK"
		be.Message = ""
		return be
	}

	// Timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		be.Stage = "Timeout"
		be.Message = err.Error()
		return be
	}

	// DNS
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		be.Stage = "DNS"
		be.Message = dnsErr.Error()
		return be
	}

	// TLS: CA unknown, expired, wrong host, etc
	var uaErr x509.UnknownAuthorityError
	if errors.As(err, &uaErr) {
		be.Stage = "TLS"
		be.Message = "TLS: unknown certificate authority"
		return be
	}

	// Network issues (refused, unreachable)
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		be.Stage = "Network"
		be.Message = opErr.Error()
		return be
	}

	if err.Error() == "redirect received but not followed (--no-follow)" {
		be.Stage = "INFO"
		be.Message = err.Error()
		return be
	}
	if err.Error() == "the parameter test flag is enabled, but no parameters were found in the specified url" {
		be.Stage = "INFO"
		be.Message = err.Error()
		return be
	}

	// Fallback
	be.Stage = "Unknown"
	be.Message = err.Error()
	return be
}
