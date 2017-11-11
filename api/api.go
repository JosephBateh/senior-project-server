package api

import (
	"log"

	"github.com/josephbateh/senior-project-server/authentication"
	db "github.com/josephbateh/senior-project-server/database"
	"github.com/zmb3/spotify"
)

// GetUserRecentlyPlayed gets the 50 most recently played songs
func GetUserRecentlyPlayed(userID string) []spotify.RecentlyPlayedItem {
	client, err := getClientFromID(userID)
	if err != nil {
		log.Fatal(err)
	}

	recents, err := client.PlayerRecentlyPlayed()
	if err != nil {
		log.Fatal(err)
	}

	return recents
}

func getClientFromID(userID string) (spotify.Client, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	client := authentication.GetClient(user.UserToken)
	return client, err
}
