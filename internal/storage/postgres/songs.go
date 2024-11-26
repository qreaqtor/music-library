package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

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
		`SELECT s.group_name, s.song, COALESCE(STRING_AGG(v.verse, '\n'), ''), s.releaseDate, COALESCE(s.link, '')
		FROM songs s LEFT JOIN verses v ON s.id = v.song_id
		WHERE s.group_name = $1 AND s.song = $2
		GROUP BY s.id, s.group_name, s.song, s.releaseDate, s.link;`

	err := s.db.QueryRow(query, song.Group, song.SongName).
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

func (s *SongsStorage) Search(ctx context.Context, search *domain.SongSearch) ([]*domain.Song, error) {
	songs := make([]*domain.Song, 0, search.Limit)
	var song domain.Song

	searchQuery := getSearchQuery(search)

	searchQuery.query = fmt.Sprintf(
		`SELECT s.group_name, s.song
		FROM songs s LEFT JOIN verses v ON s.id = v.song_id
		%s
		GROUP BY s.group_name, s.song
		LIMIT $%d OFFSET $%d;`,
		searchQuery.query,
		len(searchQuery.args)+1,
		len(searchQuery.args)+2,
	)
	searchQuery.args = append(searchQuery.args, search.Limit, search.Offset)

	rows, err := s.db.Query(searchQuery.query, searchQuery.args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&song.Group, &song.SongName)
		if err != nil {
			return nil, err
		}
		songs = append(songs, &song)
	}

	return songs, nil
}

func (s *SongsStorage) GetLyrics(ctx context.Context, song *domain.Song, batch *domain.Batch) ([]string, error) {
	opID := logmsg.ExtractOperationID(ctx)

	lyrics := make([]string, 0, batch.Limit)
	var songID uuid.UUID
	var verse string

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := "SELECT id FROM songs WHERE group_name = $1 AND song = $2"

	err = tx.QueryRow(query, song.Group, song.SongName).Scan(&songID)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug(err.Error(), "operation", opID)
		return nil, ErrUnknownResourse
	}
	if err != nil {
		return nil, err
	}

	query = "SELECT verse FROM verses WHERE song_id = $1 LIMIT $2 OFFSET $3;"
	rows, err := tx.Query(query, songID, batch.Limit, batch.Offset)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug("no lyrics", "operation", opID)
		return nil, ErrUnknownResourse
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&verse)
		if err != nil {
			return nil, err
		}

		lyrics = append(lyrics, verse)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return lyrics, nil
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
	opID := logmsg.ExtractOperationID(ctx)

	var songID uuid.UUID

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRow("SELECT id FROM songs WHERE group_name = $1 AND song = $2;", song.Group, song.SongName)
	err = row.Scan(&songID)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Debug(err.Error(), "operation", opID)
		return ErrUnknownResourse
	}
	if err != nil {
		return err
	}

	songQuery, err := getSongUpdateQuery(song, update.ToSongSchema())
	if err != nil {
		slog.Debug(err.Error(), "operation", opID)
	} else {
		_, err = tx.Exec(songQuery.query, songQuery.args...)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(`DELETE FROM verses WHERE song_id = $1;`, songID)
	if err != nil {
		return err
	}

	versesQuery, err := getLyricsUpdateQuery(songID, update.ToLyricsSchema())
	if err != nil {
		slog.Debug(err.Error(), "operation", opID)
	} else {
		_, err = tx.Exec(versesQuery.query, versesQuery.args...)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
