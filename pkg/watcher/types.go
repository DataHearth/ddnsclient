package watcher

import "errors"

var (
	ErrNilProvider = errors.New("provider can't be nil")
	ErrNilHTTP     = errors.New("http can't be nil")
)
