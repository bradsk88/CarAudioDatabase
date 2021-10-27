package auth

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
)

func NewGoogleLogin() *GoogleLogin {
	return &GoogleLogin{}
}

type GoogleLogin struct {
}

func (g *GoogleLogin) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	b, err := ioutil.ReadFile("/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, people.UserinfoEmailScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	url := config.AuthCodeURL(randState)
	http.Redirect(writer, request, url, http.StatusSeeOther)
}
