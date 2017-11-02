package smartplaylists

import (
	"log"

	"github.com/zmb3/spotify"

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
	playlistOneTracks, err := client.GetPlaylistTracks(user.UserID, playlists[2].ID)
	if err != nil {
		log.Fatal(err)
	}
	playlistOneTrackObjects := playlistOneTracks.Tracks

	// Get songs from playlist 2
	playlistTwoTracks, err := client.GetPlaylistTracks(user.UserID, playlists[3].ID)
	if err != nil {
		log.Fatal(err)
	}
	playlistTwoTrackObjects := playlistTwoTracks.Tracks

	var tracksToBeAdded []spotify.ID
	// Create array that assigns to tracksToBeAdded
	for _, object := range playlistOneTrackObjects {
		tracksToBeAdded = append(tracksToBeAdded, object.Track.ID)
	}
	for _, object := range playlistTwoTrackObjects {
		tracksToBeAdded = append(tracksToBeAdded, object.Track.ID)
	}

	// Put songs from 1 and 2 into playlist 3
	_, err = client.AddTracksToPlaylist(user.UserID, playlists[0].ID, tracksToBeAdded...)
	if err != nil {
		log.Fatal(err)
	}
}
