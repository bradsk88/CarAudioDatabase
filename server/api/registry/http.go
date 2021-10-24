package registry

import (
	"context"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/response/frequency"
	frequency2 "github.com/bradsk88/CarAudioDatabase/server/repo/response/frequency"
	"net/http"
)

func NewHTTP() *HTTP {
	return &HTTP{}
}

type HTTP struct {
}

func (h *HTTP) RegisterAll(mux *http.ServeMux) error {
	mux.HandleFunc("/db/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test passed"))
	})

	repo := frequency2.NewMySQLAmplitudeRepo()

	err := repo.Initialize(context.Background())
	if err != nil {
		return fmt.Errorf("initialize: %s", err.Error())
	}

	mux.Handle("/upload", frequency.NewUpload(repo))
	mux.Handle("/get", frequency.NewGet(repo))
	return nil
}
