package database

import "golang.org/x/oauth2"
import "time"

// User stores minimal user information
type User struct {
	UserID    string
	UserToken oauth2.Token
}

// SmartPlaylist contains all the information needed to execute a smart playlist
type SmartPlaylist struct {
	Hash  uint32
	Name  string `json:"name"`
	User  string `json:"user"`
	Rules []rule `json:"rules"`
}

// Play is a single instance of a song being played by a user
type Play struct {
	Hash  uint32
	User  string
	Track string
	Time  time.Time
}

type rule struct {
	Attribute string `json:"attribute"`
	Match     bool   `json:"match"`
	Value     string `json:"value"`
}

type test struct {
	Name  string
	Phone string
}
