package ovh

import "errors"

type (
	config map[string]interface{}
)

var (
	// ErrNilOvhConfig is thrown when OVH configuration is empty
	ErrNilOvhConfig = errors.New("OVH config is mandatory")
)
