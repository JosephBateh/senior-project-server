package database

import (
	"fmt"
	"log"
	"os"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var session *mgo.Session

// CreateUser in users collection
func CreateUser(userID string, userAuth string) {
	c := session.DB(os.Getenv("MLAB_DB")).C("users")
	err := c.Insert(&user{userID, userAuth})
	if err != nil {
		log.Fatal(err)
	}
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
