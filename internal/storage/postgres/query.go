package storage

import (
	"errors"
	"fmt"
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

	for _, field := range reflect.VisibleFields(reflect.TypeOf(update)) {
		value := updateVal.FieldByName(field.Name)
		columnTag := field.Tag.Get("db")

		if value.IsZero() || columnTag == "-" {
			continue
		}

		current := fmt.Sprintf(" %s = $%d,", columnTag, len(args)+1)
		q = fmt.Sprint(q, current)

		switch value.Interface().(type) {
		case string:
			args = append(args, value.Interface().(string))
		case time.Time:
			args = append(args, value.Interface().(time.Time))
		default:
			return nil, errors.ErrUnsupported
		}
	}

	if len(args) == 0 {
		return nil, ErrEmptySongUpdate
	}

	q = fmt.Sprintf("%s WHERE group_name = $%d AND song = $%d;",
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

func getSearchQuery(search *domain.SongSearch) *query {
	q := "WHERE"
	args := make([]any, 0)

	if search.ByGroup != "" {
		q = fmt.Sprintf("%s s.group_name ILIKE $%d", q, len(args)+1)
		args = append(args, "%"+search.ByGroup+"%")
	}

	if search.BySongName != "" {
		if len(args) != 0 {
			q += " AND"
		}
		q = fmt.Sprintf("%s s.song ILIKE $%d", q, len(args)+1)
		args = append(args, "%"+search.BySongName+"%")
	}

	if search.ByLyrics != "" {
		if len(args) != 0 {
			q += " AND"
		}
		q = fmt.Sprintf("%s EXISTS (SELECT 1 FROM verses v WHERE v.song_id = s.id AND v.verse ILIKE $%d", q, len(args)+1)
		args = append(args, "%"+search.ByLyrics+"%")
	}

	if search.ByLink != "" {
		if len(args) != 0 {
			q += " AND"
		}
		q = fmt.Sprintf("%s s.link ILIKE $%d", q, len(args)+1)
		args = append(args, "%"+search.ByLink+"%")
	}

	if !search.DateFrom.IsZero() {
		if len(args) != 0 {
			q += " AND"
		}
		q = fmt.Sprintf("%s s.releaseDate >= $%d", q, len(args)+1)
		args = append(args, search.DateFrom)
	}
	if !search.DateTo.IsZero() {
		if len(args) != 0 {
			q += " AND"
		}
		q = fmt.Sprintf("%s s.releaseDate <= $%d", q, len(args)+1)
		args = append(args, search.DateTo)
	}

	if len(args) == 0 {
		q = ""
	}

	return &query{
		query: q,
		args:  args,
	}
}
