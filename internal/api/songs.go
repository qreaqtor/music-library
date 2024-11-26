package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/qreaqtor/music-library/internal/domain"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
	"github.com/qreaqtor/music-library/pkg/web"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/qreaqtor/music-library/docs"
)

type service interface {
	Info(context.Context, *domain.Song) (*domain.SongInfo, error)
	Create(context.Context, *domain.Song) error
	Delete(context.Context, *domain.Song) error
	Update(context.Context, *domain.Song, *domain.SongUpdate) error
	GetLyrics(context.Context, *domain.Song, *domain.Batch) ([]string, error)
	Search(context.Context, *domain.SongSearch) ([]*domain.Song, error)
}

type SongsAPI struct {
	srv service

	valid *validator.Validate
}

func NewSongsAPI(srv service) *SongsAPI {
	return &SongsAPI{
		srv:   srv,
		valid: validator.New(validator.WithRequiredStructEnabled()),
	}
}

// @title Music-library API
// @version 1.0
// @description This is an implementation of an online song library
func (s *SongsAPI) Register(r *mux.Router) {
	groupAndSong := []string{
		"group", "{group:.+}",
		"song", "{song:.+}",
	}

	offsetAndLimit := []string{
		"offset", `{offset:\d+}`,
		"limit", `{limit:[1-9][\d+]?}`,
	}

	r.Path("/create").HandlerFunc(s.create).Methods(http.MethodPost)

	r.Path("/info").HandlerFunc(s.info).Methods(http.MethodGet).
		Queries(groupAndSong...)

	r.Path("/update").HandlerFunc(s.update).Methods(http.MethodPatch).
		Queries(groupAndSong...)

	r.Path("/delete").HandlerFunc(s.delete).Methods(http.MethodDelete).
		Queries(groupAndSong...)

	r.Path("/lyrics").HandlerFunc(s.getLyrics).Methods(http.MethodGet).
		Queries(append(groupAndSong, offsetAndLimit...)...)

	r.Path("/search").HandlerFunc(s.search).Methods(http.MethodGet)

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:50055/v1/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
}

// @Summary Get song info
// @Description Retrieve detailed information about a song
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} domain.SongInfo
// @Router /info [get]
func (s *SongsAPI) info(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	songInfo, err := s.srv.Info(r.Context(), song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		songInfo,
	)
}

// @Summary Get song lyrics
// @Description Retrieve lyrics of a song in batches
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Param offset query int false "Offset for batch"
// @Param limit query int false "Limit for batch"
// @Success 200 {object} map[string]any
// @Router /lyrics [get]
func (s *SongsAPI) getLyrics(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	// ignore errors, because i used regexp for this query params
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	batch := &domain.Batch{
		Offset: offset,
		Limit:  limit,
	}

	text, err := s.srv.GetLyrics(r.Context(), song, batch)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		map[string]any{
			"lyrics": text,
		},
	)
}

// @Summary Search for songs
// @Description Search for songs based on various criteria
// @Tags songs
// @Accept json
// @Produce json
// @Param search body domain.SongSearch true "Search parameters"
// @Success 200 {object} map[string]any
// @Router /search [get]
func (s *SongsAPI) search(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	search := &domain.SongSearch{}

	err := web.ReadRequestBody(r, search)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	err = s.valid.StructCtx(r.Context(), search)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	songs, err := s.srv.Search(r.Context(), search)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		map[string]any{
			"songs": songs,
		},
	)
}

// @Summary Update song information
// @Description Update details of a song including group, name, lyrics, link, and release date
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Param update body domain.SongUpdate true "Update parameters"
// @Success 200 {object} map[string]string
// @Router /update [patch]
func (s *SongsAPI) update(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	songUpdate := &domain.SongUpdate{}

	err := web.ReadRequestBody(r, songUpdate)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	err = s.valid.StructCtx(r.Context(), songUpdate)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	err = s.srv.Update(r.Context(), song, songUpdate)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusNotFound))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		map[string]string{
			"status": "ok",
		},
	)
}

// @Summary Delete a song
// @Description Remove a song from the database
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} map[string]string
// @Router /delete [delete]
func (s *SongsAPI) delete(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	err := s.srv.Delete(r.Context(), song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusNotFound))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		map[string]string{
			"status": "ok",
		},
	)
}

// @Summary Create a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body domain.Song true "Song data"
// @Success 200 {object} map[string]string
// @Router /create [post]
func (s *SongsAPI) create(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.RequestURI, r.Method)

	song := &domain.Song{}

	err := web.ReadRequestBody(r, song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
		return
	}

	err = s.valid.StructCtx(r.Context(), song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	err = s.srv.Create(r.Context(), song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusNotFound))
		return
	}

	web.WriteData(
		w,
		msg.With("OK", http.StatusOK),
		map[string]string{
			"status": "ok",
		},
	)
}
