package utils

import "errors"

var (
	// ErrReadConfigFile is thrown when viper failed to read config file
	ErrReadConfigFile = errors.New("failed to read config file")
)