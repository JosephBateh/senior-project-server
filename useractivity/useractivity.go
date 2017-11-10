package useractivity

import (
	"log"
	"time"

	db "github.com/josephbateh/senior-project-server/database"
)

var ticker *time.Ticker
var quit chan struct{}

func checkUserActivity() {
	log.Println("Begin checking user activity")
	// Get list of all users
	users, err := db.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	db.UpdatePlaysForUser(users)

	log.Println("End checking user activity")
}

// Start begins checking user activity every period (in minutes)
func Start(period int) {
	ticker = time.NewTicker(time.Duration(period) * time.Minute)
	quit := make(chan struct{})
	func() {
		for {
			select {
			case <-ticker.C:
				checkUserActivity()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Stop the tracking of user activity
func Stop() {
	ticker.Stop()
}
