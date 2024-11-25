package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/qreaqtor/music-library/internal/domain"
)

type SongsStorage struct {
	connPool *pgxpool.Pool
}

func NewSongsStorage(connPool *pgxpool.Pool) *SongsStorage {
	return &SongsStorage{
		connPool: connPool,
	}
}

func (s *SongsStorage) Get(ctx context.Context) ([]*domain.Song, error) {
	query := fmt.Sprintf("SELECT group, song FROM %s")
}
