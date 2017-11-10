package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/josephbateh/senior-project-server/smartplaylists"
	"github.com/josephbateh/senior-project-server/useractivity"

	"github.com/josephbateh/senior-project-server/authentication"
)

// Start the server
func Start() {
	setupRoutes()

	// Start the server
	var wg sync.WaitGroup
	wg.Add(1)
	go http.ListenAndServe(":8080", nil)

	// Start auto-updating playlists every N minutes
	wg.Add(2)
	go smartplaylists.Start(1)
	go useractivity.Start(1)
	log.Println("Server started")

	// Wait until go routines run
	wg.Wait()
}

func setupRoutes() {
	http.HandleFunc("/", authentication.Login)
	http.HandleFunc("/callback", authentication.Complete)
	http.HandleFunc("/smartplaylist", smartplaylists.Playlists)
}
