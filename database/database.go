package database

import (
	"fmt"
	"hash/fnv"
	"log"
	"os"

	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

// UpdatePlaysForUser updates the number of times a user has played songs
func UpdatePlaysForUser(users []User) {

}

// AddUser in users collection
func AddUser(userID string, userAuth oauth2.Token) {
	connect()
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	err := c.Insert(&User{userID, userAuth})
	if err != nil {
		log.Fatal(err)
	}
	disconnect()
}

// GetUser returns the user with the given userID
func GetUser(userID string) (User, error) {
	connect()
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	// Query One
	result := User{}
	err := c.Find(bson.M{"userid": userID}).One(&result)
	disconnect()
	return result, err
}

// GetAllUsers returns all users in the database
func GetAllUsers() ([]User, error) {
	connect()
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	// Query One
	result := []User{}
	err := c.Find(bson.M{}).All(&result)
	disconnect()
	return result, err
}

// AddSmartPlaylist adds as SmartPlaylist to the database
func AddSmartPlaylist(playlist SmartPlaylist) {
	connect()
	hashString := playlist.User + playlist.Name
	hashVal := hash(hashString)
	playlist.Hash = hashVal

	c := session.DB(os.Getenv("MLAB_DB")).C("smartplaylists")

	// Check if playlist already exists
	exists := SmartPlaylist{}
	err := c.Find(bson.M{"hash": hashVal}).One(&exists)
	if err == nil {
		// It already exists, delete it
		c.Remove(bson.M{"hash": hashVal})
	}

	err = c.Insert(playlist)
	if err != nil {
		fmt.Println(err)
	}
	disconnect()
	log.Println("Smart playlist added")
}

// GetAllSmartPlaylists returns all smart playlists in the DB
func GetAllSmartPlaylists() ([]SmartPlaylist, error) {
	connect()
	c := session.DB(os.Getenv("MLAB_DB")).C("smartplaylists")
	// Query One
	result := []SmartPlaylist{}
	err := c.Find(bson.M{}).All(&result)
	disconnect()
	return result, err
}

// Test the database connection
func Test() {
	connect()
	c := session.DB("spm-test").C("test")
	err := c.Insert(&test{"Ale", "+55 53 8116 9639"},
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
	disconnect()
}

func connect() *mgo.Session {
	var err error
	session, err = mgo.Dial(os.Getenv("MLAB_LOGIN"))
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

func disconnect() {
	session.Close()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
