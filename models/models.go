package models

import "github.com/graph-gophers/graphql-go"

type Link struct {
	Base
	Description string
	URL         string
	// this is for gorm's association
	PostedByID graphql.ID
	// this is for the schema, populated only for resolver return value
	PostedBy *User `gorm:"foreignkey:PostedByID"`
}

type User struct {
	Base
	Email    string
	Password string
	Name     string
	Links    []Link `gorm:"foreignkey:PostedByID"`
}
