package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/josephbateh/senior-project-server/smartplaylists"

	"github.com/josephbateh/senior-project-server/authentication"
)

// Start the server
func Start() {
	setupRoutes()

	// Start the server
	var wg sync.WaitGroup
	wg.Add(1)
	go http.ListenAndServe(":8080", nil)
	fmt.Println("Server listening")

	// Test it
	go smartplaylists.PlaylistFromOtherPlaylists("jbspotifytest01", "Smart Playlist", "53lV2g8Jn3cGfXtT6adA3i", "2vh7lBRsthuZMR93BrIGLX")

	// Wait until go routines run
	wg.Wait()
}

func setupRoutes() {
	http.HandleFunc("/", authentication.Login)
	http.HandleFunc("/callback", authentication.Complete)
	//http.HandleFunc("/smartplaylist", smartplaylists.PlaylistFromOtherPlaylists)
}
