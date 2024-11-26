package storage

import "errors"

var (
	ErrUnknownResourse   = errors.New("unknown resource")
	ErrEmptySongUpdate   = errors.New("nothing to update in song")
	ErrEmptyLyricsUpdate = errors.New("empty lyrics")
)
