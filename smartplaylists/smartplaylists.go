package smartplaylists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	db "github.com/josephbateh/senior-project-server/database"
	"github.com/josephbateh/senior-project-server/rest"
)

// Playlists is the function called for the smartplaylist endpoint
func Playlists(response http.ResponseWriter, request *http.Request) {
	var smartplaylist db.SmartPlaylist

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

		db.AddSmartPlaylist(smartplaylist)

		// Allow playlist to execute without blocking response
		go executeSmartPlaylist(smartplaylist)
	}

	rest.PostRequest(response, request, smartplaylist)
}

// Attributes returns the possible attributes for rules
func Attributes(response http.ResponseWriter, request *http.Request) {
	type attr struct {
		Attributes []string
	}

	butes := []string{"plays", "playlist", "artist"}

	values := attr{
		Attributes: butes,
	}

	rest.GetRequest(response, request, values)
}

func executeSmartPlaylist(smartplaylist db.SmartPlaylist) {
	userID := smartplaylist.User

	// Get the results of each rule
	tracks, err := getTracksFromRules(smartplaylist)
	if err != nil {
		return
	}

	// If playlist doesn't exist, create it
	playlistID, err := getPlaylistIDFromName(userID, smartplaylist.Name)
	if err != nil {
		playlistID = createNewPlaylist(userID, smartplaylist.Name)
	}

	updatePlaylist(userID, playlistID, tracks)
}
