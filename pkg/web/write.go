package web

import (
	"encoding/json"
	"net/http"

	logmsg "github.com/qreaqtor/music-library/pkg/logging/message"
)

/*
Вызывается в случае появления ошибки, пишет msg в логи.
Статус ответа и сообщение достает из msg.
Возвращает 500 в случае неудачной записи в w.
*/
func WriteError(w http.ResponseWriter, msg *logmsg.LogMsg) {
	msg.Error()
	http.Error(w, msg.Text, msg.Status)
}

/*
Выполняет сериализацию data и пишет в w.
В случаае появления ошибки вызывает writeError().
*/
func WriteData(w http.ResponseWriter, msg *logmsg.LogMsg, data any) {
	response, err := json.Marshal(data)
	if err != nil {
		WriteError(w, msg.With(err.Error(), http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	_, err = w.Write(response)
	if err != nil {
		WriteError(w, msg.With(err.Error(), http.StatusInternalServerError))
		return
	}

	msg.Info()
}
