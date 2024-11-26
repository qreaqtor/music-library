package domain

import "time"

type Song struct {
	Group    string `json:"group" validate:"required,min=1"`
	SongName string `json:"song" validate:"required,min=1"`
}

type SongInfo struct {
	Group       string    `json:"group"`
	SongName    string    `json:"song"`
	Lyrics      string    `json:"lyrics"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}

type SongSearch struct {
	Batch

	Search string `json:"search" valid:"required"`

	ByGroup    bool `json:"by_group" valid:"required"`
	BySongName bool `json:"by_song_name" valid:"required"`
	ByLyrics   bool `json:"by_lyrics" valid:"required"`
	ByLink     bool `json:"by_link" valid:"required"`
}
