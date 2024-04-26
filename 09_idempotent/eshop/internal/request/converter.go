package request

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"io"
	"log/slog"
	"net/http"
)

func HttpReqToObj(w http.ResponseWriter, r *http.Request, m proto.Message) bool {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("read body request", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	err = json.Unmarshal(b, m)
	if err != nil {
		slog.Error("unmarshal proto", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	return true
}
