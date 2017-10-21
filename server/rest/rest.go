package rest

// Example:
// http://localhost:4000/dummydata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Start starts the server
func Start() {
	http.HandleFunc("/playlists/0000", top10)
	http.HandleFunc("/playlists/0001", current)
	http.HandleFunc("/playlists/0002", work)
	http.HandleFunc("/playlists", playlists)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

type song struct {
	Title  string
	Artist string
	Album  string
}

func top10(writer http.ResponseWriter, response *http.Request) {
	finalResponse := [10]song{}
	finalResponse[0] = song{"It Was Always You", "Maroon 5", "V"}
	finalResponse[1] = song{"Walk Away From The Sun", "Seether", "Finding Beauty in Negative Spaces"}
	finalResponse[2] = song{"Periscope (feat. Skylar Grey)", "Papa Roach", "Crooked Teeth"}
	finalResponse[3] = song{"Faily Local", "Twenty One Pilots", "Blurryface"}
	finalResponse[4] = song{"Animals", "Maroon 5", "V"}
	finalResponse[5] = song{"Like Suicide", "Seether", "Finding Beauty in Negative Spaces"}
	finalResponse[6] = song{"Incomplete", "Backstreet Boys", "Never Gone"}
	finalResponse[7] = song{"Show Me the Meaning of Being Lonely", "Backstreet Boys", "Millenium"}
	finalResponse[8] = song{"Long Way Down", "Robert DeLong", "In The Cards"}
	finalResponse[9] = song{"Get Out Alive", "Three Days Grace", "One-X"}

	getRequest(writer, response, finalResponse)
}

func current(writer http.ResponseWriter, response *http.Request) {
	finalResponse := [5]song{}
	finalResponse[0] = song{"Ghosts 'n' Stuff", "deadmau5", "For Lack of a Better Name"}
	finalResponse[1] = song{"None Of The Above", "Papa Roach", "Crooked Teeth"}
	finalResponse[2] = song{"Periscope (feat. Skylar Grey)", "Papa Roach", "Crooked Teeth"}
	finalResponse[3] = song{"Heavy (feat. Kiiara)", "Linkin Park", "Heavy (feat. Kiiara)"}
	finalResponse[4] = song{"Your Song", "Rita Ora", "Your Song"}

	getRequest(writer, response, finalResponse)
}

func work(writer http.ResponseWriter, response *http.Request) {
	finalResponse := [3]song{}
	finalResponse[0] = song{"The Last of Us", "Gustavo Santaolalla", "The Last of Us"}
	finalResponse[1] = song{"Extraction Point", "Hans Zimmer", "Call of Duty: Modern Warfare 2"}
	finalResponse[2] = song{"Elixir", "Sizzle Bird", "Elixir"}

	getRequest(writer, response, finalResponse)
}

func playlists(writer http.ResponseWriter, response *http.Request) {
	type playlist struct {
		Title   string
		Creator string
		UUID    string
	}

	// Declare array of songs
	finalResponse := [3]playlist{}
	finalResponse[0] = playlist{"Top 10", "Joseph Bateh", "0000"}
	finalResponse[1] = playlist{"Current", "Joseph Bateh", "0001"}
	finalResponse[2] = playlist{"Work", "Joseph Bateh", "0002"}

	getRequest(writer, response, finalResponse)
}

func dummydata(writer http.ResponseWriter, response *http.Request) {
	title := "Glish"
	artist := "deadmau5"
	album := "W:/2016ALBUM/"

	// Declare array of songs
	finalResponse := [20]song{}

	// First song is different than the other 19
	finalResponse[0] = song{"Support", "deadmau5", "stuff i used to do"}
	for i := 1; i < 20; i++ {
		finalResponse[i] = song{title, artist, album}
	}

	getRequest(writer, response, finalResponse)
}

func getRequest(writer http.ResponseWriter, response *http.Request, v interface{}) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	writer.Header().Set("Content-Type", "application/json")

	// Check if the method is a get
	if response.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		fmt.Println(writer)
		return
	}

	b, err := json.Marshal(v)
	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
	}

	writer.Write(b)
}
