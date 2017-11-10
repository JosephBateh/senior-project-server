package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetRequest is used for an HTTP GET request
func GetRequest(writer http.ResponseWriter, response *http.Request, v interface{}) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	writer.Header().Set("Content-Type", "application/json")

	// Check if the method is a get
	if response.Method != http.MethodGet {
		http.Error(writer, http.StatusText(405), 405)
		fmt.Println(writer)
		return
	}

	enc := json.NewEncoder(writer)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
}

// PostRequest is used for an HTTP POST request
func PostRequest(response http.ResponseWriter, request *http.Request, v interface{}) {
	response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	response.Header().Set("Access-Control-Allow-Origin", "*")
	response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
	response.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(response)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
}
