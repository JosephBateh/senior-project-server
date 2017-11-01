package smartplaylists

import (
	"fmt"
	"log"

	"github.com/josephbateh/senior-project-server/authentication"

	"github.com/josephbateh/senior-project-server/database"
)

// Test does something
func Test() {
	firstPlaylist()
}

func firstPlaylist() {
	// Get user ID
	userID := "jbspotifytest01"
	database.Connect()

	// Get user from DB
	user, err := database.GetUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get client from user
	auth := authentication.GetAuthenticator()
	userToken := user.UserToken
	client := auth.NewClient(&userToken)

	// Get playlists
	playlistsPage, err := client.GetPlaylistsForUser(user.UserID)
	if err != nil {
		log.Fatal(err)
	}
	playlists := playlistsPage.Playlists

	// Get songs from playlist 1
	playlistOneTracks, err := client.GetPlaylistTracks(user.UserID, playlists[1].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(playlistOneTracks.Tracks[0].Track.ID)

	// Get songs from playlist 2
	playlistTwoTracks, err := client.GetPlaylistTracks(user.UserID, playlists[2].ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(playlistTwoTracks.Tracks[0].Track.ID)

	// Check if playlist exists

	// Put songs from 1 and 2 into playlist 3
	client.CreatePlaylistForUser(user.UserID, "Smart Playlist", true)

	// Disconnect for DB
	database.Disconnect()
}
