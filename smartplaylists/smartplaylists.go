package smartplaylists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zmb3/spotify"

	"github.com/josephbateh/senior-project-server/authentication"
	db "github.com/josephbateh/senior-project-server/database"
	"github.com/josephbateh/senior-project-server/rest"
)

type rule struct {
	User      string `json:"user"`
	Attribute string `json:"attribute"`
	Match     bool   `json:"match"`
	Value     string `json:"value"`
}

// Playlists is the function called for the smartplaylist endpoint
func Playlists(response http.ResponseWriter, request *http.Request) {
	var rules []rule
	if request.Method == http.MethodPost {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
		}

		err = json.Unmarshal(body, &rules)
		if err != nil {
			fmt.Println(err)
		}

		var tracks []string
		var userID string
		for i := 0; i < len(rules); i++ {
			rule := rules[i]
			ruleTracks := PlaylistMatchValue(rule.User, rule.Match, rule.Value)
			tracks = ruleTracks
			userID = rule.User
		}

		updatePlaylist(userID, tracks)
	}

	rest.PostRequest(response, request, rules)
}

func updatePlaylist(userID string, tracks []string) {
	user, client, err := getUserClient("jbspotifytest01")
	if err != nil {
		log.Fatal(err)
	}

	// Get all track IDs in one slice
	var trackIDs []spotify.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, spotify.ID(track))
	}

	// Clear playlist and add tracks
	tracksCurrentlyInPlaylist, _ := client.GetPlaylistTracks(user.UserID, spotify.ID("0EVAoDQ4B8Rsd2rtiE9AyO"))
	var currentTrackIDs []spotify.ID
	for _, object := range tracksCurrentlyInPlaylist.Tracks {
		currentTrackIDs = append(currentTrackIDs, object.Track.ID)
	}

	client.RemoveTracksFromPlaylist(user.UserID, spotify.ID("0EVAoDQ4B8Rsd2rtiE9AyO"), currentTrackIDs...)
	client.AddTracksToPlaylist(user.UserID, spotify.ID("0EVAoDQ4B8Rsd2rtiE9AyO"), trackIDs...)
	log.Println("Playlist updated")
}

func getUserClient(userID string) (db.User, spotify.Client, error) {
	// Get user from the DB
	user, err := db.GetUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get client from user
	client := authentication.GetClient(user.UserToken)
	return user, client, err
}

// PlaylistMatchValue will return tracks that are in the provided playlist
func PlaylistMatchValue(userID string, match bool, value string) []string {
	user, client, err := getUserClient(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get users playlists
	playlistPage, err := client.GetPlaylist(user.UserID, spotify.ID(value))
	if err != nil {
		log.Fatal(err)
	}
	playlistTracks := playlistPage.Tracks.Tracks

	var tracks []string
	for i := 0; i < len(playlistTracks); i++ {
		track := playlistTracks[i].Track.ID
		tracks = append(tracks, string(track))
	}
	return tracks
}
