package smartplaylists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/josephbateh/senior-project-server/rest"
)

type rule struct {
	Attribute string `json:"attribute"`
	Match     bool   `json:"match"`
	Value     string `json:"value"`
}

type smartplaylist struct {
	Name  string `json:"name"`
	User  string `json:"user"`
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

		// Get the user ID
		userID := smartplaylist.User

		// Get the results of each rule
		tracks := getTracksFromRules(smartplaylist)

		// If playlist doesn't exist, create it
		playlistID, err := getPlaylistIDFromName(userID, smartplaylist.Name)
		if err != nil {
			playlistID = createNewPlaylist(userID, smartplaylist.Name)
		}

		// Clear playlist and add new tracks
		updatePlaylist(userID, playlistID, tracks)
	}

	rest.PostRequest(response, request, smartplaylist)
}

func runSmartPlaylist(smartplaylist smartplaylist) {

}
