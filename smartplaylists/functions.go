package smartplaylists

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/josephbateh/senior-project-server/authentication"
	db "github.com/josephbateh/senior-project-server/database"
	"github.com/zmb3/spotify"
	set "gopkg.in/fatih/set.v0"
)

type ruleFunc func(string, string, string) ([]string, error)

const (
	plays        = "plays"
	playlist     = "playlist"
	artistString = "artist"
)

const (
	equal    = "is"
	notEqual = "isnot"
	greater  = "greater"
	less     = "less"
)

func ruleFunctions() map[string]ruleFunc {
	var ruleFunctions = map[string]ruleFunc{
		playlist:     playlistFunc,
		plays:        playsFunc,
		artistString: artistFunc,
	}
	return ruleFunctions
}

func getTracksFromRules(smartplaylist db.SmartPlaylist) ([]string, error) {
	ruleFunctions := ruleFunctions()

	var isMatch [][]string
	var isNotMatch [][]string
	var isGreaterMatch [][]string
	var isLessMatch [][]string

	for i := 0; i < len(smartplaylist.Rules); i++ {
		rule := smartplaylist.Rules[i]
		ruleTracks, err := ruleFunctions[rule.Attribute](smartplaylist.User, rule.Match, rule.Value)
		if err != nil {
			return nil, err
		}
		switch rule.Match {
		case equal:
			isMatch = append(isMatch, ruleTracks)
		case notEqual:
			isNotMatch = append(isNotMatch, ruleTracks)
		case greater:
			isGreaterMatch = append(isGreaterMatch, ruleTracks)
		case less:
			isLessMatch = append(isLessMatch, ruleTracks)
		}
	}

	union := unionOfTracks(isMatch...)
	union = append(union, unionOfTracks(isGreaterMatch...)...)
	union = append(union, unionOfTracks(isLessMatch...)...)
	intersection := unionOfTracks(isNotMatch...)

	return intersectionOfTracks(union, intersection), nil
}

func getUserClient(userID string) (db.User, spotify.Client, error) {
	// Get user from the DB
	user, err := db.GetUser(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get client from user
	client := authentication.GetClient(user.UserToken)
	return user, client, err
}

func artistFunc(userID string, match string, value string) ([]string, error) {
	_, client, err := getUserClient(userID)
	if err != nil {
		fmt.Println(err)
	}

	searchResult, err := client.Search(value, spotify.SearchTypeArtist)
	var artistID spotify.ID
	if err == nil {
		artistID = searchResult.Artists.Artists[0].ID
	} else {
		return nil, err
	}

	// Get artist albums
	artist, err := client.GetArtistAlbums(artistID)
	if err != nil {
		return nil, err
	}
	artistAlbums := artist.Albums

	var tracks []string
	for _, album := range artistAlbums {
		albumID := album.ID
		albumTracksPage, err := client.GetAlbumTracks(albumID)
		if err == nil {
			albumTracks := albumTracksPage.Tracks
			for _, track := range albumTracks {
				tracks = append(tracks, track.ID.String())
			}
		}
	}

	return tracks, err
}

func playlistFunc(userID string, match string, value string) ([]string, error) {
	user, client, err := getUserClient(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get users playlists
	playlistPage, err := client.GetPlaylist(user.UserID, spotify.ID(value))
	if err != nil {
		fmt.Println(err)
	}
	playlistTracks := playlistPage.Tracks.Tracks

	var tracks []string
	for i := 0; i < len(playlistTracks); i++ {
		track := playlistTracks[i].Track.ID
		tracks = append(tracks, string(track))
	}
	return tracks, nil
}

func playsFunc(userID string, match string, value string) ([]string, error) {
	plays, err := db.NumberOfPlays(userID)
	if err != nil {
		return nil, err
	}

	var playedSongs []string
	var playSet = set.New()
	for _, play := range plays {
		playSet.Add(play.Track)
	}

	playedSongs = set.StringSlice(playSet)

	var result []string
	for _, play := range playedSongs {
		numPlays := db.NumberOfPlaysForTrack(userID, play)
		valueInt, _ := strconv.Atoi(value)
		switch match {
		case "is":
			if numPlays == valueInt {
				result = append(result, play)
			}
		case "is not":
			if numPlays != valueInt {
				result = append(result, play)
			}
		case "greater":
			if numPlays > valueInt {
				result = append(result, play)
			}
		case "less":
			if numPlays < valueInt {
				result = append(result, play)
			}
		default:
			fmt.Println("default")
		}
	}
	return result, nil
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
		log.Println("Error getting playlist", err)
		return "", err
	}

	simplePlaylistPage, err := client.GetPlaylistsForUser(userID)
	var simplePlaylistArray []spotify.SimplePlaylist
	if err != nil {
		log.Println("Error getting playlist", err)
		simplePlaylistArray = nil
	} else {
		simplePlaylistArray = simplePlaylistPage.Playlists
	}

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
		log.Println("Error getting client", err)
		return
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

	if len(currentTrackIDs) > 100 {
		div := len(currentTrackIDs) / 100
		mod := len(currentTrackIDs) % 100
		for i := 0; div > 0; div-- {
			rge := i * 100
			i++
			var currIDs []spotify.ID
			if div == 0 {
				currIDs = currentTrackIDs[rge : rge+mod]
			} else {
				currIDs = currentTrackIDs[rge : rge+99]
			}
			client.RemoveTracksFromPlaylist(user.UserID, playlistID, currIDs...)
		}
	} else {
		client.RemoveTracksFromPlaylist(user.UserID, playlistID, currentTrackIDs...)
	}

	if len(trackIDs) > 100 {
		times := len(trackIDs) / 100
		mod := len(trackIDs) % 100
		for i := 0; times > 0; times-- {
			rge := i * 100
			i++
			var currIDs []spotify.ID
			if times == 0 {
				currIDs = trackIDs[rge : rge+mod]
			} else {
				currIDs = trackIDs[rge : rge+99]
			}
			client.AddTracksToPlaylist(user.UserID, playlistID, currIDs...)
		}
	} else {
		client.AddTracksToPlaylist(user.UserID, playlistID, trackIDs...)
	}
}

func createNewPlaylist(userID string, name string) string {
	_, client, err := getUserClient(userID)
	if err != nil {
		log.Println("Error getting client", err)
		return ""
	}

	playlist, err := client.CreatePlaylistForUser(userID, name, false)
	if err != nil {
		log.Println("Error getting playlist", err)
		return ""
	}

	return string(playlist.ID)
}
