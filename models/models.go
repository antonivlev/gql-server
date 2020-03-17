package models

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

type Link struct {
	ID          graphql.ID `gorm:"primary_key"`
	Description string
	URL         string
}

// BeforeCreate will set a UUID rather than numeric ID.
func (link *Link) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid.String())
}
