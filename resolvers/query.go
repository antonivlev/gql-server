package resolvers

import (
	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
)

func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil
}

func (r *RootResolver) Feed() ([]models.Link, error) {
	return database.GetAllLinks()
}
