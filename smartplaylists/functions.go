package smartplaylists

import (
	"errors"
	"log"

	"github.com/zmb3/spotify"
	set "gopkg.in/fatih/set.v0"
)

func getTracksFromRules(smartplaylist smartplaylist) []string {
	var trueMatch [][]string
	var falseMatch [][]string

	for i := 0; i < len(smartplaylist.Rules); i++ {
		rule := smartplaylist.Rules[i]
		ruleTracks := playlistMatchValue(rule.User, rule.Match, rule.Value)
		if rule.Match {
			trueMatch = append(trueMatch, ruleTracks)
		} else {
			falseMatch = append(falseMatch, ruleTracks)
		}
	}

	unionOfTrue := unionOfTracks(trueMatch...)
	unionOfFalse := unionOfTracks(falseMatch...)

	return intersectionOfTracks(unionOfTrue, unionOfFalse)
}

func getUserIDFromSmartPlaylist(smartplaylist smartplaylist) string {
	return smartplaylist.Rules[0].User
}

func unionOfTracks(trackList ...[]string) []string {
	tracks := set.New()

	// TODO: Make this not O(N^2)
	for i := 0; i < len(trackList); i++ {
		newSet := set.New()
		for j := 0; j < len(trackList[i]); j++ {
			newSet.Add(trackList[i][j])
		}
		tracks.Merge(newSet)
	}

	return set.StringSlice(tracks)
}

// This function has not been tested yet
func intersectionOfTracks(original []string, trackList ...[]string) []string {
	tracks := set.New()

	for i := 0; i < len(original); i++ {
		tracks.Add(original[i])
	}

	// TODO: Make this not O(N^2)
	for i := 0; i < len(trackList); i++ {
		newSet := set.New()
		for j := 0; j < len(trackList[i]); j++ {
			newSet.Add(trackList[i][j])
		}
		tracks.Separate(newSet)
	}

	return set.StringSlice(tracks)
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

// PlaylistMatchValue will return tracks that are in the provided playlist
func playlistMatchValue(userID string, match bool, value string) []string {
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
