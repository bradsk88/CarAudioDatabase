package frequency

import (
	"context"
	"encoding/json"
	"github.com/bradsk88/CarAudioDatabase/server/api/common"
	model "github.com/bradsk88/CarAudioDatabase/server/model/frequency"
	"log"
	"net/http"
)

type Getter interface {
	Get(
		ctx context.Context, id string,
	) ([]model.DataPoint, error)
}

func NewGet(
	getter Getter,
) *Get {
	return &Get{
		getter: getter,
	}
}

type Get struct {
	getter Getter
}

type GetResponse struct {
	Data []model.DataPoint `json:"data"`
}

func (g *Get) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Serving %s\n", request.URL)

	common.EnableCors(writer)
	id := "7c2d9cd3-6b6b-412f-be19-8f8e0d57e4cc" // TODO: Get from request
	data, err := g.getter.Get(request.Context(), id)
	if err != nil {
		log.Printf("Create: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}

	out, err := json.Marshal(GetResponse{
		Data: data,
	})
	if err != nil {
		log.Printf("json.Marshal: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}

	_, err = writer.Write(out)
	if err != nil {
		log.Printf("Write: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}
}
