package rest

// Example:
// http://localhost:3000/numberconverter?number=500000&oldBase=10&newBase=16

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../converter"
)

// Start starts the server
func Start() {
	//muxRouter := http.NewServeMux()
	//muxRouter.HandleFunc("/numberconverter", numberconverter)
	//http.Handle("/", muxRouter)
	http.HandleFunc("/numberconverter", numberconverter)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

type reqres struct {
	Number  string
	OldBase string
	NewBase string
	Result  string
}

func numberconverter(writer http.ResponseWriter, response *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
	writer.Header().Set("Content-Type", "application/json")

	// Check if the method is a get
	if response.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		fmt.Println(writer)
		return
	}

	number := response.FormValue("number")
	oldBase := response.FormValue("oldBase")
	newBase := response.FormValue("newBase")
	result := converter.ConvertStringNumberToNewBase(number, oldBase, newBase)
	finalResponse := reqres{number, oldBase, newBase, result}
	b, err := json.Marshal(finalResponse)
	if err != nil {
		http.Error(writer, http.StatusText(500), 500)
	}
	writer.Write(b)
}
