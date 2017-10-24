package authentication

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth      spotify.Authenticator
	waitGroup sync.WaitGroup
	ch        = make(chan *spotify.Client)
	state     = "abc123"
)

// Start starts the process of listening for authentication requests
func Start() {

	loadEnv()

	waitGroup.Add(1)
	go listen()
	fmt.Println("Searver listening...")
	waitGroup.Wait()

}

func listen() {
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryModify, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserTopRead)

	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", login)

	log.Fatal(http.ListenAndServe(":8080", nil))
	waitGroup.Done()
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login...")
	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}
	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	ch <- &client
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
