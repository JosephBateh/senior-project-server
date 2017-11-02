package authentication

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"

	"github.com/josephbateh/senior-project-server/database"
	"github.com/josephbateh/senior-project-server/rest"
	"github.com/zmb3/spotify"
)

var (
	waitGroup sync.WaitGroup
	ch        = make(chan *spotify.Client)
	state     = "u4KEsvUyfQ9O"
)

func redirectURI() string {
	return os.Getenv("REDIRECT_URI")
}

func getAuthenticator() spotify.Authenticator {
	authenticator := spotify.NewAuthenticator(redirectURI(), spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryModify, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserTopRead)
	return authenticator
}

// GetClient returns a client for making API calls
func GetClient(token oauth2.Token) spotify.Client {
	auth := getAuthenticator()
	return auth.NewClient(&token)
}

// Login to the server
func Login(w http.ResponseWriter, r *http.Request) {
	auth := getAuthenticator()
	loginURL := auth.AuthURL(state)

	type res struct {
		Address string
	}
	response := res{loginURL}
	rest.GetRequest(w, r, response)
}

// Complete login to the server
func Complete(w http.ResponseWriter, r *http.Request) {
	auth := getAuthenticator()
	// Get token
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}

	// Exchange the token for a new client
	client := auth.NewClient(tok)

	// Let the user know that server auth has completed
	fmt.Fprintf(w, "Login Completed! You may now close this tab.")

	// Get the userID
	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	// Get the client's token for storage
	clientToken, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}

	// Check if a user with that ID already exists
	// If they do not exist, add them
	_, err = database.GetUser(user.ID)
	if err != nil {
		database.AddUser(user.ID, *clientToken)
	}
	log.Println("User -", user.ID, "logged in.")
}
