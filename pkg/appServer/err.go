package appserver

import "errors"

var (
	ErrAlreadyStarted = errors.New("server already started")
	ErrNotStarted     = errors.New("server not started")
)
