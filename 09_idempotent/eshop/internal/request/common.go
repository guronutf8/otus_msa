package request

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type Common struct {
	Status  bool     `json:"status"`
	Message string   `json:"message,omitempty"`
	Log     []string `json:"log,omitempty"`
}

func (c Common) WriteResponse(w http.ResponseWriter) {
	indent, err := json.MarshalIndent(c, "", "	")
	if err != nil {
		slog.Error("common to []byte", "err", err)
		return
	}
	if _, err := w.Write(indent); err != nil {
		slog.Error("common write response", "err", err)
	}
}

func FromBody(data io.ReadCloser) (*Common, error) {
	all, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	r := Common{}
	if err := json.Unmarshal(all, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
