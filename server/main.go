package main

import (
	"bufio"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/registry"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

func main() {

	mux := http.NewServeMux()

	ipFileName := "./car-audio-database/dist/car-audio-database/ipAddress.js"
	f, err := os.OpenFile(ipFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}

	ipAddress := os.Getenv("CARAV_DB_IP")
	if ipAddress == "" {
		ipAddress = "localhost"
	}

	_, err = f.Write([]byte(fmt.Sprintf("document.carAVServerIP = \"%s\";", ipAddress)))
	if err != nil {
		log.Fatalf("ipAddress.Write: %s", err.Error())
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./car-audio-database/dist/car-audio-database"))
	mux.Handle("/", fs)

	sessionKey, err := getSessionKey()
	if err != nil {
		log.Fatalf("getSessionKey: %s", err.Error())
		return
	}

	store := sessions.NewCookieStore(sessionKey)

	reg := registry.NewHTTP(ipAddress)
	err = reg.RegisterAll(mux, store)
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

func getSessionKey() ([]byte, error) {
	f, err := os.Open("/session.key")
	if err != nil {
		return nil, fmt.Errorf("os.Open: %s", err.Error())
	}

	rd := bufio.NewReader(f)
	key, _, err := rd.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("ReadLine: %s", err.Error())
	}

	return key, nil
}
