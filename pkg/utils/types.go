package utils

import "errors"

var (
	// ErrReadConfigFile is thrown when viper failed to read config file
	ErrReadConfigFile = errors.New("failed to read config file")
	// ErrNilLogger is thrown when the parameter logger is nil
	ErrNilLogger = errors.New("logger can't be nil")
	// ErrWrongStatusCode is thrown when the response status code isn't a 200
	ErrWrongStatusCode = errors.New("response sent an non 200 status code")
)
