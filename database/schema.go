package database

import "golang.org/x/oauth2"

type test struct {
	Name  string
	Phone string
}

type user struct {
	UserID    string
	UserToken oauth2.Token
}
