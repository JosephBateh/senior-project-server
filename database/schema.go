package database

import "golang.org/x/oauth2"

type test struct {
	Name  string
	Phone string
}

// User stores minimal user information
type User struct {
	UserID    string
	UserToken oauth2.Token
}
