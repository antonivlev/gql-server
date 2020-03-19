package resolvers

import (
	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
)

func (r *RootResolver) Post(args struct {
	Description string
	URL         string
}) (models.Link, error) {
	newLink := models.Link{
		URL:         args.URL,
		Description: args.Description,
	}

	dbLink, errCreate := database.CreateLink(newLink)
	if errCreate != nil {
		return models.Link{}, errCreate
	}
	return *dbLink, nil
}

// Note: if field is nullable in schema, corresponding field in struct must be pointer (so it can be nil)
type AuthPayload struct {
	Token *string
	User  *models.User
}

func (r *RootResolver) Signup(args struct {
	Email    string
	Password string
	Name     string
}) (*AuthPayload, error) {
	u, errCreate := database.CreateUser(args.Email, args.Password, args.Name)
	if errCreate != nil {
		return nil, errCreate
	}

	payload := &AuthPayload{
		Token: nil,
		User:  u,
	}
	return payload, nil
}

func (r *RootResolver) Login(args struct {
	Email    string
	Password string
}) (*AuthPayload, error) {
	return nil, nil
}
