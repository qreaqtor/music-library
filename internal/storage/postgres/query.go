package storage

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/qreaqtor/music-library/internal/domain"
)

type query struct {
	query string
	args  []any
}

// using reflect for getting tag for column name in database.
// support only string and time.Time types for update.SongSchema fields
func getSongUpdateQuery(song *domain.Song, update domain.SongSchema) (*query, error) {
	q := "UPDATE songs SET"
	args := make([]any, 0)

	updateVal := reflect.ValueOf(update)

	for i, field := range reflect.VisibleFields(reflect.TypeOf(update)) {
		value := updateVal.FieldByName(field.Name)
		columnTag := field.Tag.Get("db")

		if value.IsZero() || columnTag == "-" {
			continue
		}

		current := fmt.Sprintf(" %s = $%d,", columnTag, i+1)
		q = fmt.Sprint(q, current)

		switch value.Interface().(type) {
		case string:
			args = append(args, value.String())
		case time.Time:
			args = append(args, value.Interface().(time.Time))
		default:
			return nil, errors.ErrUnsupported
		}
	}

	slog.Debug(fmt.Sprint(args))

	if len(args) == 0 {
		return nil, ErrEmptySongUpdate
	}

	q = fmt.Sprintf("%s WHERE group_name = $%d AND song = $%d RETURNING id;",
		q[:len(q)-1],
		len(args)+1,
		len(args)+2,
	)
	args = append(args, song.Group, song.SongName)

	return &query{
		query: q,
		args:  args,
	}, nil
}

func getLyricsUpdateQuery(songID uuid.UUID, update domain.LyricsSchema) (*query, error) {
	if len(update.Lyrics) == 0 {
		return nil, ErrEmptyLyricsUpdate
	}

	numbers := make([]string, 0, len(update.Lyrics))
	args := make([]any, 0, len(update.Lyrics)+1)
	args = append(args, songID)

	for i, verse := range update.Lyrics {
		numbers = append(numbers, fmt.Sprintf("($1, $%d)", i+2))
		args = append(args, verse)
	}

	q := fmt.Sprintf(
		`INSERT INTO verses (song_id, verse) VALUES %s`,
		strings.Join(numbers, ","),
	)

	return &query{
		query: q,
		args:  args,
	}, nil
}
