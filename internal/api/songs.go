package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qreaqtor/music-library/internal/domain"
	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
	"github.com/qreaqtor/music-library/pkg/web"
)

type service interface {
	Info(*domain.Song) (*domain.SongInfo, error)
	Create(*domain.Song) error
	Delete(*domain.Song) error
	Update(*domain.Song, *domain.SongUpdate) error
}

type SongsAPI struct {
	srv service
}

func NewSongsAPI(srv service) *SongsAPI {
	return &SongsAPI{
		srv: srv,
	}
}

func (s *SongsAPI) Register(r mux.Router) {
	r.Path("/search").HandlerFunc(s.search).Methods(http.MethodGet)

	r.Path("/create").HandlerFunc(s.create).Methods(http.MethodPost)

	rQueries := r.NewRoute().Subrouter().Queries(
		"group", "{group:.+}",
		"song", "{song:.+}",
	)

	rQueries.Path("/info").HandlerFunc(s.info).Methods(http.MethodGet)

	rQueries.Path("/update").HandlerFunc(s.update).Methods(http.MethodPatch)

	rQueries.Path("/delete").HandlerFunc(s.delete).Methods(http.MethodDelete)
}

func (s *SongsAPI) info(w http.ResponseWriter, r *http.Request) {
	msg := logmsg.NewLogMsg(r.Context(), r.URL.Path, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	songInfo, err := s.srv.Info(song)
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

	err = web.ValidateStruct(r.Context(), songUpdate)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	err = s.srv.Update(song, songUpdate)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
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
	msg := logmsg.NewLogMsg(r.Context(), r.URL.Path, r.Method)

	song := &domain.Song{
		Group:    r.URL.Query().Get("group"),
		SongName: r.URL.Query().Get("song"),
	}

	err := s.srv.Delete(song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
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

	err = web.ValidateStruct(r.Context(), song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusUnprocessableEntity))
		return
	}

	err = s.srv.Create(song)
	if err != nil {
		web.WriteError(w, msg.With(err.Error(), http.StatusBadRequest))
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
