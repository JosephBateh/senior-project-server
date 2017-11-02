package server

import (
	"fmt"
	"sync"

	"github.com/josephbateh/senior-project-server/smartplaylists"

	"github.com/josephbateh/senior-project-server/authentication"
)

// Start the server
func Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go authentication.Listen()
	fmt.Println("Server listening")
	go smartplaylists.PlaylistFromOtherPlaylists("jbspotifytest01", "Smart Playlist", "1gnsYyxX6gNEgmkVpeGNTK", "2vh7lBRsthuZMR93BrIGLX")
	wg.Wait()
}
