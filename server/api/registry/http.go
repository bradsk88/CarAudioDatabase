package registry

import (
	"context"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/auth"
	"github.com/bradsk88/CarAudioDatabase/server/api/response/frequency"
	frequency2 "github.com/bradsk88/CarAudioDatabase/server/repo/response/frequency"
	"github.com/bradsk88/CarAudioDatabase/server/repo/users"
	"github.com/gorilla/sessions"
	"net/http"
)

func NewHTTP(ipAddr string) *HTTP {
	return &HTTP{
		ipAddr: ipAddr,
	}
}

type HTTP struct {
	ipAddr string
}

func (h *HTTP) RegisterAll(mux *http.ServeMux, sess *sessions.CookieStore) error {
	frRepo := frequency2.NewMySQLAmplitudeRepo()

	err := frRepo.Initialize(context.Background())
	if err != nil {
		return fmt.Errorf("frRepo.initialize: %s", err.Error())
	}

	userRepo := users.NewRepo()

	err = userRepo.Initialize(context.Background())
	if err != nil {
		return fmt.Errorf("userRepo.initialize: %s", err.Error())
	}

	mux.Handle("/upload", frequency.NewUpload(frRepo, sess))
	mux.Handle("/get", frequency.NewGet(frRepo))
	mux.Handle("/google-callback", auth.NewGoogleCallback(sess, userRepo))
	mux.Handle("/google-login", auth.NewGoogleLogin())
	return nil
}
