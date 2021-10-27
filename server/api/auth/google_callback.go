package auth

import (
	"context"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/keys"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
	"io/ioutil"
	"log"
	"net/http"
)

type Ensurer interface {
	Ensure(
		ctx context.Context, googleId string, googleEmailAddress string,
	) (userId string, err error)
}

func NewGoogleCallback(
	sess *sessions.CookieStore, ensurer Ensurer,
) *GoogleCallback {
	return &GoogleCallback{
		sess:    sess,
		ensurer: ensurer,
	}
}

type GoogleCallback struct {
	sess    *sessions.CookieStore
	ensurer Ensurer
}

func (g *GoogleCallback) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	session, _ := g.sess.Get(request, keys.SessionName)

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
	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Printf("config.Exchange: %s\n", err.Error())
		writer.WriteHeader(400)
		_, _ = writer.Write([]byte("code not accepted"))
		return
	}
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

	userId, err := g.ensurer.Ensure(ctx, r.Id, r.Email)
	if err != nil {
		log.Printf("Failed to ensure user %v", err)
		writer.WriteHeader(500)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = userId

	err = session.Save(request, writer)
	if err != nil {
		log.Printf("Failed to save session %v", err)
		writer.WriteHeader(500)
		return
	}

	_, _ = writer.Write([]byte(fmt.Sprintf("Logged in as %s", userId)))
}
