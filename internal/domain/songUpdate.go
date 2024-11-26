package domain

import "time"

type SongUpdate struct {
	Group       string    `json:"group" validate:"omitempty,min=1"`
	SongName    string    `json:"song" validate:"omitempty,min=1"`
	Lyrics      []string  `json:"lyrics" validate:"omitempty,min=1"`
	Link        string    `json:"link" validate:"omitempty,http_url"`
	ReleaseDate time.Time `json:"releaseDate" validate:"omitempty"`
}

type SongSchema struct {
	Group       string    `db:"group_name"`
	SongName    string    `db:"song"`
	Link        string    `db:"link"`
	ReleaseDate time.Time `db:"releaseDate"`
}

type LyricsSchema struct {
	Lyrics []string
}

func (s *SongUpdate) ToSongSchema() SongSchema {
	return SongSchema{
		Group:       s.Group,
		SongName:    s.SongName,
		Link:        s.Link,
		ReleaseDate: s.ReleaseDate,
	}
}

func (s *SongUpdate) ToLyricsSchema() LyricsSchema {
	return LyricsSchema{
		Lyrics: s.Lyrics,
	}
}
