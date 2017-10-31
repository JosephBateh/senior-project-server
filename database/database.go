package database

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

// AddUser in users collection
func AddUser(userID string, userAuth oauth2.Token) {
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	err := c.Insert(&User{userID, userAuth})
	if err != nil {
		log.Fatal(err)
	}
}

// GetUser returns the user with the given userID
func GetUser(userID string) (User, error) {
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	// Query One
	result := User{}
	err := c.Find(bson.M{"userid": userID}).One(&result)
	return result, err
}

// Connect to the database
func Connect() *mgo.Session {
	var err error
	session, err = mgo.Dial(os.Getenv("MLAB_LOGIN"))
	if err != nil {
		panic(err)
	}
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	return session
}

// Disconnect from the database
func Disconnect() {
	session.Close()
}

// Test the database connection
func Test() {
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
}
