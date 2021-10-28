package frequency

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/common"
	model "github.com/bradsk88/CarAudioDatabase/server/model/frequency"
	"github.com/gorilla/sessions"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type Creator interface {
	Create(
		ctx context.Context, createdByUserId string, data []byte,
	) error
}

func NewUpload(creator Creator, sess *sessions.CookieStore) *Upload {
	return &Upload{
		creator: creator,
	}
}

type Upload struct {
	creator Creator
}

func (u *Upload) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	common.EnableCors(writer)

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := request.ParseMultipartForm(10 << 20)
	if err != nil {
		log.Printf("req.ParseMultipartForm: %s", err)
		writer.WriteHeader(500)
	}

	file, _, err := request.FormFile("file")
	if err != nil {
		log.Printf("request.FormFile: %s", err)
		return
	}
	defer func() {
		err := file.Close()
		if err != nil {
			writer.WriteHeader(500)
			log.Printf("file.Close: %s", err)
		}
	}()

	fr, err := captureData(file)
	if err != nil {
		writer.WriteHeader(400)
		_, err = writer.Write([]byte(fmt.Sprintf("Could not extract data: %s", err.Error())))
		if err != nil {
			log.Println(err)
		}
	}

	res, err := json.Marshal(fr)
	if err != nil {
		writer.WriteHeader(500)
		return
	}

	err = u.creator.Create(request.Context(), "bradsk88", res)
	if err != nil {
		log.Printf("Create: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}
}

func captureData(file multipart.File) ([]model.DataPoint, error) {
	scanner := bufio.NewScanner(file)

	fr := make([]model.DataPoint, 0, 20000*5)
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

func parseLine(line string) (*model.DataPoint, error) {
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
	return &model.DataPoint{Frequency: freq, Amplitude: amp, Phase: phase}, nil
}
