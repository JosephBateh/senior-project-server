package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

const redirectURI := DEV_BASE_URL + "callback/"

var (
	auth  = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	baseURL := os.Getenv("DEV_BASE_URL")
	fmt.Println(baseURL)

	redirectURL := baseURL + "callback/"
	fmt.Println(redirectURL)

	authenticate(redirectURL)
}

func authenticate(redirectURL string) {
	auth := spotify.NewAuthenticator(redirectURL, spotify.ScopePlaylistReadPrivate)

	auth.SetAuthInfo(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"))

	url := auth.AuthURL("state-string")
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token("state-string", r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	// create a client using the specified token
	client := auth.NewClient(token)

	// the client can now be used to make authenticated requests
}
