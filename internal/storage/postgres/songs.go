package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/qreaqtor/music-library/internal/domain"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
)

type SongsStorage struct {
	db *sql.DB
}

func NewSongsStorage(connection *sql.DB) *SongsStorage {
	return &SongsStorage{
		db: connection,
	}
}

func (s *SongsStorage) Info(ctx context.Context, song *domain.Song) (*domain.SongInfo, error) {
	opID := logmsg.ExtractOperationID(ctx)

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
		slog.Debug(err.Error(), "operation", opID)
		return nil, ErrUnknownResourse
	}
	if err != nil {
		return nil, err
	}

	return songInfo, nil
}

func (s *SongsStorage) Create(ctx context.Context, song *domain.Song) error {
	query := "INSERT INTO songs (group_name, song) values ($1, $2);"
	_, err := s.db.Exec(query, song.Group, song.SongName)
	if err != nil {
		return err
	}
	return nil
}

func (s *SongsStorage) Delete(ctx context.Context, song *domain.Song) error {
	query := "DELETE FROM songs WHERE group_name = $1 AND song = $2;"

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
	var id uuid.UUID

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query :=
		`UPDATE songs SET group_name = $1, song = $2, releaseDate = $3, link = $4
		WHERE group_name = $5 AND song = $6
		RETURNING id;`

	row := tx.QueryRow(
		query,
		update.Group,
		update.SongName,
		update.ReleaseDate,
		update.Link,
		song.Group,
		song.SongName,
	)
	err = row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug(err.Error(), "operation", logmsg.ExtractOperationID(ctx))
		return ErrUnknownResourse
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM verses WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	numbers := make([]string, 0, len(update.Lyrics))
	args := make([]any, 0, len(update.Lyrics)+1)
	args = append(args, id)
	for i, verse := range update.Lyrics {
		numbers = append(numbers, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, verse)
	}

	query = fmt.Sprintf(
		`INSERT INTO verses (id, verse) VALUES %s`,
		strings.Join(numbers, ","),
	)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
