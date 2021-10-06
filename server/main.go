package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./car-audio-database/dist/car-audio-database"))
	mux.Handle("/", fs)

	mux.HandleFunc("/endpoint", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("test passed"))
	})

	mux.HandleFunc("/upload", uploadFile)

	fmt.Println("Serving...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("failed to start server: %s\n", err.Error())
	}
}
