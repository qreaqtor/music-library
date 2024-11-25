package service

import "github.com/qreaqtor/music-library/internal/domain"

type storage interface {
	Info(*domain.Song) (*domain.SongInfo, error)
	Create(*domain.Song) error
	Delete(*domain.Song) error
	Update(*domain.Song, *domain.SongUpdate) error
}

type SongsService struct {
	st storage
}

func (s *SongsService) Info(song *domain.Song) (*domain.SongInfo, error) {
	return s.st.Info(song)
}

func (s *SongsService) Create(song *domain.Song) error {
	return s.st.Create(song)
}

func (s *SongsService) Delete(song *domain.Song) error {
	return s.st.Delete(song)
}

func (s *SongsService) Update(song *domain.Song, update *domain.SongUpdate) error {
	return s.st.Update(song, update)
}
