package database

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"

	"github.com/zmb3/spotify"

	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UpdatePlaysForUser updates the number of times a user has played songs
func UpdatePlaysForUser(user User, recents []spotify.RecentlyPlayedItem) {
	for _, track := range recents {
		err := addSongPlay(user.UserID, track)
		if err != nil {
			//log.Println("Failed to add play:", err)
		}
	}
}

func addSongPlay(userID string, track spotify.RecentlyPlayedItem) error {
	session, err := connect()
	c := session.DB(os.Getenv("MLAB_DB")).C("plays")
	if err != nil {
		return err
	}

	// Get hash for play
	hash := hash(userID + track.PlayedAt.String())

	// Check if play already exists
	result := Play{}
	err = c.Find(bson.M{"hash": hash, "user": userID}).One(&result)
	if err != nil {
		// Play not found
		err = c.Insert(&Play{hash, userID, string(track.Track.ID), track.PlayedAt})
		if err != nil {
			log.Fatal(err)
		}
	}
	disconnect(session)
	return err
}

// NumberOfPlaysForTrack returns the number of times a user has played a track
func NumberOfPlaysForTrack(user string, track string) (int, error) {
	session, err := connect()
	if err != nil {
		return 0, err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("plays")

	// Check if play already exists
	results := []Play{}
	err = c.Find(bson.M{"user": user, "track": track}).All(&results)
	plays := 0
	if err == nil {
		// Played at least once
		plays = len(results)
	}
	disconnect(session)
	return plays, nil
}

// NumberOfPlays returns the plays for a user
func NumberOfPlays(user string) ([]Play, error) {
	session, err := connect()
	if err != nil {
		return nil, err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("plays")

	results := []Play{}
	err = c.Find(bson.M{}).All(&results)
	if err != nil {
		return nil, err
	}
	disconnect(session)
	return results, nil
}

// AddUser in users collection
func AddUser(userID string, userAuth oauth2.Token) error {
	session, err := connect()
	if err != nil {
		return err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	err = c.Insert(&User{userID, userAuth})
	if err != nil {
		log.Fatal(err)
	}
	disconnect(session)
	return err
}

// GetUser returns the user with the given userID
func GetUser(userID string) (User, error) {
	session, err := connect()
	result := User{}
	if err != nil {
		return result, err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	// Query One
	err = c.Find(bson.M{"userid": userID}).One(&result)
	disconnect(session)
	return result, err
}

// GetAllUsers returns all users in the database
func GetAllUsers() ([]User, error) {
	session, err := connect()
	if err != nil {
		return nil, err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("users")

	// Query One
	result := []User{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	disconnect(session)
	return result, err
}

// AddSmartPlaylist adds as SmartPlaylist to the database
func AddSmartPlaylist(playlist SmartPlaylist) error {
	session, err := connect()
	if err != nil {
		return err
	}

	hashString := playlist.User + playlist.Name
	hashVal := hash(hashString)
	playlist.Hash = hashVal

	c := session.DB(os.Getenv("MLAB_DB")).C("smartplaylists")

	// Check if playlist already exists
	exists := SmartPlaylist{}
	err = c.Find(bson.M{"hash": hashVal}).One(&exists)
	if err == nil {
		// It already exists, delete it
		c.Remove(bson.M{"hash": hashVal})
	}

	err = c.Insert(playlist)
	if err != nil {
		fmt.Println(err)
	}
	disconnect(session)
	log.Println("Smart playlist added")
	return err
}

// GetAllSmartPlaylists returns all smart playlists in the DB
func GetAllSmartPlaylists() ([]SmartPlaylist, error) {
	session, err := connect()
	if err != nil {
		return nil, err
	}
	c := session.DB(os.Getenv("MLAB_DB")).C("smartplaylists")
	// Query One
	result := []SmartPlaylist{}
	err = c.Find(bson.M{}).All(&result)
	if err != nil {
		return nil, err
	}
	disconnect(session)
	return result, err
}

// Test the database connection
func Test() error {
	session, err := connect()
	if err != nil {
		return err
	}
	c := session.DB("spm-test").C("test")
	err = c.Insert(&test{"Ale", "+55 53 8116 9639"},
		&test{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := test{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
	disconnect(session)
	return err
}

func connect() (*mgo.Session, error) {
	session, err := mgo.Dial(os.Getenv("MLAB_LOGIN"))
	if err == nil {
		session.SetMode(mgo.Monotonic, true)
	}
	return session, err
}

func disconnect(session *mgo.Session) {
	session.Close()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
