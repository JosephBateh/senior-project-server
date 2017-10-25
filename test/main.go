package main

import (
	"fmt"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type person struct {
	Name  string
	Phone string
}

func main() {
	fmt.Println("Attempting to dial!")
	session, err := mgo.Dial("mongodb://admin:Pa55w0rd@ds229415.mlab.com:29415/spm-test")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("spm-test").C("people")
	err = c.Insert(&person{"Ale", "+55 53 8116 9639"},
		&person{"Cla", "+55 53 8402 8510"})
	if err != nil {
		log.Fatal(err)
	}

	result := person{}
	err = c.Find(bson.M{"name": "Ale"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Phone:", result.Phone)
}
