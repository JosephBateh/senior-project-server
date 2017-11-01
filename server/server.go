package server

import (
	"fmt"
	"sync"

	"github.com/josephbateh/senior-project-server/authentication"
)

// Start the server
func Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go authentication.Listen()
	fmt.Println("Server listening")
	wg.Wait()
}
