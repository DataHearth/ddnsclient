package utils

import (
	"errors"
)

// * Errors
var (
	// ErrReadConfigFile is thrown when viper failed to read config file
	ErrReadConfigFile = errors.New("failed to read config file")
	// ErrNilLogger is thrown when the parameter logger is nil
	ErrNilLogger = errors.New("logger is mandatory")
	// ErrNilOvhConfig is thrown when OVH configuration is empty
	ErrNilOvhConfig = errors.New("OVH config is mandatory")
	// ErrNilProvider ...
	ErrNilProvider = errors.New("provider is mandatory")
	// ErrNilHTTP ...
	ErrNilHTTP = errors.New("http is mandatory")
	// ErrWrongStatusCode is thrown when the response status code isn't a 200
	ErrWrongStatusCode = errors.New("response sent an non 200 status code")
	// ErrGetServerIP is thrown when HTTP can't contact the web-ip service
	ErrGetServerIP = errors.New("failed to fetch server IP")
	// ErrParseHTTPBody is thrown when the HTTP service can't parse the body response
	ErrParseHTTPBody = errors.New("can't parse response body")
	// ErrHeadRemoteIP ...
	ErrHeadRemoteIP = errors.New("failed to fetch subdomain informations with HEAD")
	// ErrSplitAddr ...
	ErrSplitAddr = errors.New("can't split subdomain remote IP address")
	// ErrCreateNewRequest ...
	ErrCreateNewRequest = errors.New("can't create http request")
	// ErrUpdateRequest ...
	ErrUpdateRequest = errors.New("failed to set new IP address")
)

type (
	ProviderConfig map[string]interface{}
)
