package registry

import (
	"github.com/bradsk88/CarAudioDatabase/server/api/response/frequency"
	"net/http"
)

func NewHTTP() *HTTP {
	return &HTTP{}
}

type HTTP struct {
}

func (h *HTTP) RegisterAll(mux *http.ServeMux) {
	mux.HandleFunc("/db/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test passed"))
	})

	mux.Handle("/upload", frequency.NewUpload())
}
