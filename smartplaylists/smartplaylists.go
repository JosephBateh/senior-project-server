package smartplaylists

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zmb3/spotify"

	"github.com/josephbateh/senior-project-server/authentication"
	db "github.com/josephbateh/senior-project-server/database"
	"github.com/josephbateh/senior-project-server/rest"
	"gopkg.in/fatih/set.v0"
)

type rule struct {
	User      string `json:"user"`
	Attribute string `json:"attribute"`
	Match     bool   `json:"match"`
	Value     string `json:"value"`
}

type smartplaylist struct {
	Name  string `json:"name"`
	Rules []rule `json:"rules"`
}

// Playlists is the function called for the smartplaylist endpoint
func Playlists(response http.ResponseWriter, request *http.Request) {
	var smartplaylist smartplaylist

	// If statement guards against an OPTIONS request
	if request.Method == http.MethodPost {

		// Parse JSON into byte array
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err)
		}

		// Put new smart playlist information into smartplaylist
		err = json.Unmarshal(body, &smartplaylist)
		if err != nil {
			fmt.Println(err)
		}

		// Get the results of each rule
		var tracks [][]string
		var userID string
		for i := 0; i < len(smartplaylist.Rules); i++ {
			rule := smartplaylist.Rules[i]
			ruleTracks := PlaylistMatchValue(rule.User, rule.Match, rule.Value)
			tracks = append(tracks, ruleTracks)
			userID = rule.User
		}

		// Check if playlist already exists
		playlistID, err := getPlaylistIDFromName(userID, smartplaylist.Name)

		// If playlist doesn't exist, create it
		if err != nil {
			playlistID = createNewPlaylist(userID, smartplaylist.Name)
		}

		// Clear playlist and add new tracks
		updatePlaylist(userID, playlistID, unionOfTracks(tracks...))
	}

	rest.PostRequest(response, request, smartplaylist)
}

func getPlaylistIDFromName(userID string, name string) (string, error) {
	_, client, err := getUserClient(userID)
	if err != nil {
		log.Fatal(err)
	}

	simplePlaylistPage, err := client.GetPlaylistsForUser(userID)
	if err != nil {
		log.Fatal(err)
	}

	simplePlaylistArray := simplePlaylistPage.Playlists

	var playlistID string

	for _, playlist := range simplePlaylistArray {
		playlistName := playlist.Name
		if playlistName == name {
			err = nil
			return string(playlist.ID), err
		}
	}

	err = errors.New("No playlist with that ID")

	return playlistID, err
}

func updatePlaylist(userID string, playlistIDString string, tracks []string) {
	playlistID := spotify.ID(playlistIDString)

	user, client, err := getUserClient(userID)
	if err != nil {
		log.Fatal(err)
	}

	// Get all track IDs in one slice
	var trackIDs []spotify.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, spotify.ID(track))
	}

	// Clear playlist and add tracks
	tracksCurrentlyInPlaylist, _ := client.GetPlaylistTracks(user.UserID, playlistID)
	var currentTrackIDs []spotify.ID
	for _, object := range tracksCurrentlyInPlaylist.Tracks {
		currentTrackIDs = append(currentTrackIDs, object.Track.ID)
	}

	client.RemoveTracksFromPlaylist(user.UserID, playlistID, currentTrackIDs...)
	client.AddTracksToPlaylist(user.UserID, playlistID, trackIDs...)
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

// TODO: Make this not O(N^2)
func unionOfTracks(trackList ...[]string) []string {
	tracks := set.New()

	for i := 0; i < len(trackList); i++ {
		newSet := set.New()
		for j := 0; j < len(trackList[i]); j++ {
			newSet.Add(trackList[i][j])
		}
		tracks.Merge(newSet)
	}

	return set.StringSlice(tracks)
}

func createNewPlaylist(userID string, name string) string {
	_, client, err := getUserClient(userID)
	if err != nil {
		log.Fatal(err)
	}

	playlist, err := client.CreatePlaylistForUser(userID, name, false)
	if err != nil {
		log.Fatal(err)
	}

	return string(playlist.ID)
}
