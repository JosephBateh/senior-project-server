package server

import (
	"fmt"
	"sync"

	"github.com/josephbateh/senior-project-server/authentication"
	"github.com/josephbateh/senior-project-server/smartplaylists"
)

// Start the server
func Start() {
	var wg sync.WaitGroup
	wg.Add(2)
	go authentication.Listen()
	go smartplaylists.Test()
	fmt.Println("Server listening")
	wg.Wait()
}
