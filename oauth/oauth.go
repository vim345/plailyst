package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

// LocalOauth holds info about LocalOauth object.
type LocalOauth struct {
	Config *oauth2.Config
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

// Login logs in a user.
func (o *LocalOauth) Login(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	u := o.Config.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// GetUserDataFromGoogle uses code to get token and get user info from Google.
func (o *LocalOauth) GetService(code string) (*youtube.Service, error) {
	ctx := context.Background()
	token, err := o.Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	youtubeService, err := youtube.NewService(ctx, option.WithTokenSource(o.Config.TokenSource(ctx, token)))
	if err != nil {
		return nil, fmt.Errorf("Got error creating the service: %s", err.Error())
	}
	return youtubeService, nil
}

// NewOauth creates a new Oauth Object.
func NewOauth() *LocalOauth {
	secret := os.Getenv("CLIENT_SECRET")
	if len(secret) == 0 {
		secret = "client_secret.json"
	}
	b, err := ioutil.ReadFile(secret)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return &LocalOauth{Config: config}
}
