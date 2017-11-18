package useractivity

import (
	"log"
	"time"

	"github.com/josephbateh/senior-project-server/api"

	db "github.com/josephbateh/senior-project-server/database"
)

var ticker *time.Ticker
var quit chan struct{}

func checkUserActivity() {
	log.Println("Begin checking user activity")
	start := time.Now()
	// Get list of all users
	users, err := db.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range users {
		recents := api.GetUserRecentlyPlayed(user.UserID)
		db.UpdatePlaysForUser(user, recents)
	}

	elapsed := time.Since(start)
	userCount := len(users)
	log.Printf("Updated %v users in %s...", userCount, elapsed)
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
