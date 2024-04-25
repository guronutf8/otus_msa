package delivery

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Slot    int32  `json:"slot"`
}

type Rollback struct {
	Slot int32 `json:"slot"`
}

func (r Response) WriteResponse(w http.ResponseWriter) {
	indent, err := json.MarshalIndent(r, "	", "")
	if err != nil {
		slog.Error("common to []byte", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(indent); err != nil {
		slog.Error("common write response", "err", err)
	}
}

func FromBody(data io.ReadCloser) (*Response, error) {
	all, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	r := Response{}
	if err := json.Unmarshal(all, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
