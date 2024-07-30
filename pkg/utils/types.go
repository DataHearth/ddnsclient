package utils

import (
	"errors"
)

// * Errors
var (
	// ErrNilLogger is thrown when the parameter logger is nil
	ErrNilLogger = errors.New("logger is mandatory")
	// ErrWrongStatusCode is thrown when the response status code isn't a 200
	ErrWrongStatusCode = errors.New("web-ip returns a non 200 status code")
	// ErrGetServerIP is thrown when HTTP can't contact the web-ip service
	ErrGetServerIP = errors.New("HTTP error")
	// ErrParseHTTPBody is thrown when the HTTP service can't parse the body response
	ErrParseHTTPBody = errors.New("can't parse response body")
	// ErrSplitAddr is thrown when the remote address can't be splitted
	ErrSplitAddr = errors.New("can't split subdomain remote IP address")
	// ErrCreateNewRequest is thrown when http request creation failed
	ErrCreateNewRequest = errors.New("can't create http request")
	// ErrUpdateRequest is thrown when the update request failed
	ErrUpdateRequest = errors.New("failed to set new IP address")
	// ErrInvalidURL is thrown when user does not provide a URL and it does not exist in default urls
	ErrInvalidURL = errors.New("no url was provided")
	// ErrInvalidName is thrown when provider name was not provided
	ErrInvalidName = errors.New("no provider name was provided")
	// ErrReadBody is thrown when body response can't be parsed
	ErrReadBody = errors.New("failed to read response body")
	// ErrNilWatcher is thrown when no watcher config was provided
	ErrNilWatcher = errors.New("watcher is mandatory")
	// ErrIpLength is thrown when subdomain no or multiples remote IP address
	ErrIpLenght = errors.New("zero or more than 1 ips have been found")
	// ErrNilConfig is thrown when an empty config is provided
	ErrNilConfig = errors.New("config is mandatory")
)
