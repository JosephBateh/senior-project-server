package smartplaylists

import (
	"log"
	"time"

	db "github.com/josephbateh/senior-project-server/database"
)

var ticker *time.Ticker
var quit chan struct{}

// UpdateSmartPlaylists is public so I can test it
func UpdateSmartPlaylists() {
	start := time.Now()
	allSmartPlaylists, err := db.GetAllSmartPlaylists()
	if err != nil {
		log.Println("Error getting smart playlists:", err)
		return
	}

	for _, playlist := range allSmartPlaylists {
		executeSmartPlaylist(playlist)
	}
	elapsed := time.Since(start)
	playlists := len(allSmartPlaylists)
	log.Printf("Updated %v smart playlists in %s...", playlists, elapsed)
}

// Start updating all smart playlists every period (in minutes)
func Start(period int) {
	ticker = time.NewTicker(time.Duration(period) * time.Minute)
	quit := make(chan struct{})
	func() {
		for {
			select {
			case <-ticker.C:
				UpdateSmartPlaylists()
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
