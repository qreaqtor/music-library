package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/qreaqtor/music-library/internal/domain"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
)

type SongsStorage struct {
	db *sql.DB
}

func NewSongsStorage(connPool *sql.DB) *SongsStorage {
	return &SongsStorage{
		db: connPool,
	}
}

func (s *SongsStorage) Info(ctx context.Context, song *domain.Song) (*domain.SongInfo, error) {
	songInfo := &domain.SongInfo{}

	query :=
		`SELECT s.group_name, s.song, COALESCE(STRING_AGG(v.verse, '\n'), '') AS lyrics, s.releaseDate, COALESCE(s.link, '')
FROM songs s LEFT JOIN verses v ON s.id = v.id
WHERE s.group_name = $1 AND s.song = $2
GROUP BY s.id, s.group_name, s.song, s.releaseDate, s.link;`

	err := s.db.
		QueryRow(
			query,
			song.Group,
			song.SongName,
		).
		Scan(
			&songInfo.Group,
			&songInfo.SongName,
			&songInfo.Lyrics,
			&songInfo.ReleaseDate,
			&songInfo.Link,
		)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug(err.Error(), "operation", logmsg.ExtractOperationID(ctx))
		return nil, ErrUnknownResourse
	}
	if err != nil {
		return nil, err
	}

	return songInfo, nil
}

func (s *SongsStorage) Create(ctx context.Context, song *domain.Song) error {
	query := "INSERT INTO songs (group_name, song) values ($1, $2)"
	_, err := s.db.Exec(query, song.Group, song.SongName)
	if err != nil {
		return err
	}
	return nil
}

func (s *SongsStorage) Delete(ctx context.Context, song *domain.Song) error {
	query := "DELETE FROM songs WHERE group_name = $1 AND song = $2"
	res, err := s.db.Exec(query, song.Group, song.SongName)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil || n == 0 {
		slog.Debug("no rows affected", "operation", logmsg.ExtractOperationID(ctx))
		return ErrUnknownResourse
	}
	return nil
}

func (s *SongsStorage) Update(ctx context.Context, song *domain.Song, update *domain.SongUpdate) error {
	return nil
}
