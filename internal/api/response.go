package api

import "github.com/qreaqtor/music-library/internal/domain"

type getLyricsResponse struct {
	Lyrics []string
}

type searchResponse struct {
	Songs []*domain.Song
}

type messageResponse struct {
	Message string
}
