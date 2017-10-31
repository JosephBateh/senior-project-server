package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/joho/godotenv"
	"github.com/josephbateh/senior-project-server/database"
	"github.com/zmb3/spotify"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth      spotify.Authenticator
	waitGroup sync.WaitGroup
	ch        = make(chan *spotify.Client)
	state     = "u4KEsvUyfQ9O"
)

// Start starts the process of listening for authentication requests
func Start() {

	loadEnv()

	database.Connect()

	waitGroup.Add(1)
	go listen()
	fmt.Println("Server listening...")
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
	loginURL := auth.AuthURL(state)

	type res struct {
		Address string
	}
	response := res{loginURL}
	getRequest(w, r, response)
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
	fmt.Fprintf(w, "Login Completed! You may now close this tab.")

	// use the client to make calls that require authorization
	user, err := client.CurrentUser()
	clientToken, err := client.Token()
	if err != nil {
		log.Fatal(err)
	}
	_, err = database.GetUser(user.ID)
	if err == nil {
		database.AddUser(user.ID, *clientToken)
	}

	database.Disconnect()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("User", user.ID, "logged in")
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

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
