package models

import (
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/jinzhu/gorm"
)

// Base contains fields present in all records: ID
type Base struct {
	ID graphql.ID `gorm:"primary_key"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid.String())
}
