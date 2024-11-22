package domain

import "time"

type Song struct {
	Group       string    `json:"group"`
	SongName    string    `json:"song"`
	Lyrics      []string  `json:"lyrics"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}
