package displayname

import (
	"context"
	"fmt"
	"github.com/bradsk88/CarAudioDatabase/server/api/common"
	"github.com/bradsk88/CarAudioDatabase/server/keys"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
)

type Claimer interface {
	ClaimDisplayName(
		ctx context.Context, userID string, displayName string,
	) error
}

func NewClaim(
	sess *sessions.CookieStore, claimer Claimer,
) *GoogleCallback {
	return &GoogleCallback{
		sess:    sess,
		claimer: claimer,
	}
}

type GoogleCallback struct {
	sess    *sessions.CookieStore
	claimer Claimer
}

func (g *GoogleCallback) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	fmt.Printf("Serving %s\n", request.URL.Path)

	common.EnableCors(writer)

	// TODO: Extract "get session user ID" to a reusable service
	session, err := g.sess.Get(request, keys.SessionName)
	if err != nil {
		writer.WriteHeader(500)
		log.Printf("sess.Get: %s\n", err.Error())
		return
	}

	_userID, ok := session.Values[keys.SessionKeyUserID]
	if !ok {
		writer.WriteHeader(401)
		return
	}
	userID, ok := _userID.(string)
	if !ok {
		writer.WriteHeader(500)
		log.Printf("userID not string")
		return
	}

	dn := request.URL.Query().Get("displayName")
	if dn == "" {
		writer.WriteHeader(400)
		_, _ = writer.Write([]byte("displayName is required"))
		return
	}

	err = g.claimer.ClaimDisplayName(request.Context(), userID, dn)
	if err != nil {
		log.Printf("ClaimDisplayName: %s\n", err.Error())
		writer.WriteHeader(500)
		return
	}
}

