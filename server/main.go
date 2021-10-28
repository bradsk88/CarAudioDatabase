package main

import (
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/registry"
	"net/http"

	"github.com/gorilla/sessions"
)

func main() {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./car-audio-database/dist/car-audio-database"))
	mux.Handle("/", fs)

	key := []byte("super-secret-key") // TODO: Generate and store
	store := sessions.NewCookieStore(key)

	reg := registry.NewHTTP()
	err := reg.RegisterAll(mux, store)
	if err != nil {
		fmt.Printf("reg.RegisterAll: %s\n", err.Error())
		return
	}

	fmt.Println("Serving...")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err.Error())
	}
}
