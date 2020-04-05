package models

import "github.com/graph-gophers/graphql-go"

type Link struct {
	Base
	Description string
	URL         string
	// this is for gorm's association
	PostedByID graphql.ID
	// this should be filled automatically by gorm based on PostedByID
	PostedBy *User
}

type User struct {
	Base
	Email    string
	Password string
	Name     string
	Links    []Link `gorm:"foreignkey:PostedByID"`
}
