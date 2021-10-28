package auth

import (
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"io/ioutil"
	"log"
	"net/http"
)

func NewGoogleCallback(sess *sessions.CookieStore) *GoogleCallback {
	return &GoogleCallback{
		sess: sess,
	}
}

type GoogleCallback struct {
	sess *sessions.CookieStore
}

func (g *GoogleCallback) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	session, _ := g.sess.Get(request, "caravdb-session")

	b, err := ioutil.ReadFile("/credentials.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		writer.WriteHeader(500)
		return
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, people.UserinfoEmailScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		writer.WriteHeader(500)
		return
	}

	code := request.URL.Query().Get("code")
	ctx := request.Context()
	token, _ := config.Exchange(ctx, code) // TODO: Handle err
	client := config.Client(ctx, token)

	srv, err := oauth2.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to create oauth2 Client %v", err)
		writer.WriteHeader(500)
		return
	}

	r, err := srv.Userinfo.V2.Me.Get().Do()
	if err != nil {
		log.Printf("Unable to get user info %v", err)
		writer.WriteHeader(500)
		return
	}

	log.Printf("Login successful for %s", r.Email)
	session.Values["authenticated"] = true
	session.Values["google_id"] = r.Id

	err = session.Save(request, writer)
	if err != nil {
		log.Printf("Failed to save session %v", err)
		writer.WriteHeader(500)
		return
	}
}
