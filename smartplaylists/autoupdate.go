package smartplaylists

import (
	"fmt"
	"log"
	"time"

	db "github.com/josephbateh/senior-project-server/database"
)

var ticker *time.Ticker
var quit chan struct{}

func updateSmartPlaylists() {
	allSmartPlaylists, err := db.GetAllSmartPlaylists()
	if err != nil {
		fmt.Println(err)
	}

	for _, playlist := range allSmartPlaylists {
		executeSmartPlaylist(playlist)
	}
	log.Println("Smart Playlists Updated")
}

// Start updating all smart playlists every period (in minutes)
func Start(period int) {
	ticker = time.NewTicker(time.Duration(period) * time.Minute)
	quit := make(chan struct{})
	func() {
		for {
			select {
			case <-ticker.C:
				updateSmartPlaylists()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop the updating of smart playlists
func Stop() {
	ticker.Stop()
}
