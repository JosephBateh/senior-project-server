package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"

	"github.com/josephbateh/senior-project-server/database"
	"github.com/zmb3/spotify"
)

var (
	waitGroup   sync.WaitGroup
	ch          = make(chan *spotify.Client)
	state       = "u4KEsvUyfQ9O"
	redirectURI string
)

// Listen for authentication requests
func Listen() {
	redirectURI = os.Getenv("REDIRECT_URI")
	waitGroup.Add(1)
	go start()
	waitGroup.Wait()
}

func getAuthenticator() spotify.Authenticator {
	authenticator := spotify.NewAuthenticator(redirectURI, spotify.ScopeUserReadPrivate, spotify.ScopePlaylistReadPrivate, spotify.ScopePlaylistModifyPublic, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadCollaborative, spotify.ScopeUserLibraryModify, spotify.ScopeUserLibraryRead, spotify.ScopeUserReadPrivate, spotify.ScopeUserReadCurrentlyPlaying, spotify.ScopeUserReadPlaybackState, spotify.ScopeUserModifyPlaybackState, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserTopRead)
	return authenticator
}

// GetClient returns a client for making API calls
func GetClient(token oauth2.Token) spotify.Client {
	auth := getAuthenticator()
	return auth.NewClient(&token)
}

func start() {
	http.HandleFunc("/callback", userLogin)
	http.HandleFunc("/", login)

	log.Fatal(http.ListenAndServe(":8080", nil))
	waitGroup.Done()
}

func login(w http.ResponseWriter, r *http.Request) {
	auth := getAuthenticator()
	loginURL := auth.AuthURL(state)

	type res struct {
		Address string
	}
	response := res{loginURL}
	getRequest(w, r, response)
}

func userLogin(w http.ResponseWriter, r *http.Request) {
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

func getRequest(writer http.ResponseWriter, response *http.Request, v interface{}) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	writer.Header().Set("Content-Type", "application/json")

	// Check if the method is a get
	if response.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		fmt.Println(writer)
		return
	}

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.Encode(v)

	// b, err := json.Marshal(v)
	// if err != nil {
	// 	http.Error(writer, http.StatusText(500), 500)
	// }

	// writer.Write(b)
}
