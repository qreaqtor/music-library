package api

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/qreaqtor/music-library/internal/domain"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
	"github.com/qreaqtor/music-library/pkg/web"
)

type service interface {
	Info(context.Context, *domain.Song) (*domain.SongInfo, error)
	Create(context.Context, *domain.Song) error
	Delete(context.Context, *domain.Song) error
	Update(context.Context, *domain.Song, *domain.SongUpdate) error
}

type SongsAPI struct {
	srv service

	valid *validator.Validate
}

func NewSongsAPI(srv service) *SongsAPI {
	return &SongsAPI{
		srv: srv,
		valid: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (s *SongsAPI) Register(r *mux.Router) {
	r.Path("/search").HandlerFunc(s.search).Methods(http.MethodGet)

	r.Path("/create").HandlerFunc(s.create).Methods(http.MethodPost)

	r.Path("/info").HandlerFunc(s.info).Methods(http.MethodGet).
		Queries(
			"group", "{group:.+}",
			"song", "{song:.+}",
		)

	r.Path("/update").HandlerFunc(s.update).Methods(http.MethodPatch).
		Queries(
			"group", "{group:.+}",
			"song", "{song:.+}",
		)

	r.Path("/delete").HandlerFunc(s.delete).Methods(http.MethodDelete).
		Queries(
			"group", "{group:.+}",
			"song", "{song:.+}",
		)
}

func (s *SongsAPI) info(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.URL.Path, r.Method)

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

func (s *SongsAPI) search(w http.ResponseWriter, r *http.Request) {

}

func (s *SongsAPI) update(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.URL.Path, r.Method)

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

func (s *SongsAPI) create(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.URL.Path, r.Method)

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
