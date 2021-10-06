package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type DataPoint struct {
	Frequency float64
	Amplitude float64
	Phase     float64
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(500)
	}

	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			w.WriteHeader(500)
		}
	}()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fr, err := captureData(file)
	if err != nil {
		w.WriteHeader(400)
		_, err = w.Write([]byte(fmt.Sprintf("Could not extract data: %s", err.Error())))
		if err != nil {
			log.Fatal(err)
		}
	}

	res, err := json.MarshalIndent(fr, "", "\t")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	fmt.Println(string(res))
}

func captureData(file multipart.File) ([]DataPoint, error) {
	scanner := bufio.NewScanner(file)

	fr := make([]DataPoint, 0, 20000*5)
	i := 0
	startRead := false

	for scanner.Scan() {
		line := scanner.Text()

		if !startRead {
			if line == "* Freq(Hz) SPL(dB) Phase(degrees)" {
				startRead = true
			}
			continue
		}

		dp, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("parse line %d: %s", i+1, err.Error())
		}

		fr = append(fr, *dp)
		i++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %s", err.Error())
	}

	return fr, nil
}

func parseLine(line string) (*DataPoint, error) {
	spl := strings.Split(line, " ")
	freq, err := strconv.ParseFloat(spl[0], 64)
	if err != nil {
		return nil, fmt.Errorf("parse frequency: %s", err.Error())
	}
	amp, err := strconv.ParseFloat(spl[1], 64)
	if err != nil {
		return nil, fmt.Errorf("parse amplitude: %s", err.Error())
	}
	phase, err := strconv.ParseFloat(spl[2], 64)
	if err != nil {
		return nil, fmt.Errorf("parse phase: %s", err.Error())
	}
	return &DataPoint{Frequency: freq, Amplitude: amp, Phase: phase}, nil
}
