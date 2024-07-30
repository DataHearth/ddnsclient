package utils

import (
	"errors"
)

// * Errors
var (
	// ErrNilLogger is thrown when the parameter logger is nil
	ErrNilLogger = errors.New("logger is mandatory")
	// ErrNilProvider ...
	ErrNilProvider = errors.New("provider is mandatory")
	// ErrWrongStatusCode is thrown when the response status code isn't a 200
	ErrWrongStatusCode = errors.New("web-ip returns a non 200 status code")
	// ErrGetServerIP is thrown when HTTP can't contact the web-ip service
	ErrGetServerIP = errors.New("HTTP error")
	// ErrParseHTTPBody is thrown when the HTTP service can't parse the body response
	ErrParseHTTPBody = errors.New("can't parse response body")
	// ErrSplitAddr ...
	ErrSplitAddr = errors.New("can't split subdomain remote IP address")
	// ErrCreateNewRequest ...
	ErrCreateNewRequest = errors.New("can't create http request")
	// ErrUpdateRequest ...
	ErrUpdateRequest = errors.New("failed to set new IP address")
)
