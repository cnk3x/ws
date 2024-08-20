package ws

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Error(w http.ResponseWriter, err error, status int) {
	slog.Error("respond err", "err", err)
	Respond(w, map[string]string{"err": err.Error()}, status)
}

func Write[T ~string | ~[]byte](w http.ResponseWriter, data T, contentType string, status int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	w.Write([]byte(data))
}

func Respond(w http.ResponseWriter, obj any, status int) {
	data, err := json.Marshal(obj)
	if err != nil {
		slog.Error("respond obj error", "err", err)
		Write(w, []byte(err.Error()), "text/plain", status)
	} else {
		Write(w, data, "application/json", status)
	}
}

func Decode(r *http.Request, obj any) (err error) {
	return json.NewDecoder(r.Body).Decode(obj)
}
