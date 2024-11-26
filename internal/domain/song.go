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

	ByGroup    string `json:"by_group" valid:"omitempty,min=1"`
	BySongName string `json:"by_song_name" valid:"omitempty,min=1"`
	ByLyrics   string `json:"by_lyrics" valid:"omitempty,min=1"`
	ByLink     string `json:"by_link" valid:"omitempty,min=1"`

	DateFrom time.Time `json:"from" valid:"omitempty"`
	DateTo   time.Time `json:"to" valid:"omitempty"`
}
