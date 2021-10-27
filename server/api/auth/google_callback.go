package auth

import (
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"io/ioutil"
	"log"
	"net/http"
)

func NewGoogleCallback() *GoogleCallback {
	return &GoogleCallback{}
}

type GoogleCallback struct {
}

func (g *GoogleCallback) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	b, err := ioutil.ReadFile("/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, people.UserinfoEmailScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	code := request.URL.Query().Get("code")
	ctx := request.Context()
	token, _ := config.Exchange(ctx, code) // TODO: Handle err
	client := config.Client(ctx, token)

	srv, err := oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people Client %v", err)
	}

	r, err := srv.Userinfo.V2.Me.Get().Do()
	_, _ = writer.Write([]byte(r.Email))
}
