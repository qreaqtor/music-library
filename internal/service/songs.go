package service

import (
	"context"

	"github.com/qreaqtor/music-library/internal/domain"
)

type storage interface {
	Info(context.Context, *domain.Song) (*domain.SongInfo, error)
	Create(context.Context, *domain.Song) error
	Delete(context.Context, *domain.Song) error
	Update(context.Context, *domain.Song, *domain.SongUpdate) error
	GetLyrics(context.Context, *domain.Song, *domain.Batch) ([]string, error)
	Search(context.Context, *domain.SongSearch) ([]*domain.Song, error)
}

type SongsService struct {
	st storage
}

func NewSongsService(storage storage) *SongsService {
	return &SongsService{
		st: storage,
	}
}

func (s *SongsService) GetLyrics(ctx context.Context, song *domain.Song, batch *domain.Batch) ([]string, error) {
	return s.st.GetLyrics(ctx, song, batch)
}

func (s *SongsService) Search(ctx context.Context, search *domain.SongSearch) ([]*domain.Song, error) {
	return s.st.Search(ctx, search)
}

func (s *SongsService) Info(ctx context.Context, song *domain.Song) (*domain.SongInfo, error) {
	return s.st.Info(ctx, song)
}

func (s *SongsService) Create(ctx context.Context, song *domain.Song) error {
	return s.st.Create(ctx, song)
}

func (s *SongsService) Delete(ctx context.Context, song *domain.Song) error {
	return s.st.Delete(ctx, song)
}

func (s *SongsService) Update(ctx context.Context, song *domain.Song, update *domain.SongUpdate) error {
	return s.st.Update(ctx, song, update)
}
