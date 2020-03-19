package models

import "github.com/graph-gophers/graphql-go"

type Link struct {
	Base
	Description string
	URL         string
	PostedByID  graphql.ID
}

type User struct {
	Base
	Email    string
	Password string
	Name     string
	Links    []Link `gorm:"foreignkey:PostedByID"`
}
