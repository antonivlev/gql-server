package resolvers

import (
	"github.com/antonivlev/gql-server/auth"
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
	_, errCreate := database.CreateUser(args.Email, args.Password, args.Name)
	if errCreate != nil {
		return nil, errCreate
	}

	// todo: useless moving around from one struct to another
	var emailPassword = struct {
		Email    string
		Password string
	}{
		Email:    args.Email,
		Password: args.Password,
	}

	return r.Login(emailPassword)
}

func (r *RootResolver) Login(args struct {
	Email    string
	Password string
}) (*AuthPayload, error) {
	u, errAuth := database.GetUserByCredentials(args.Email, args.Password)
	if errAuth != nil {
		return nil, errAuth
	}

	token, errToken := auth.GenerateToken()
	if errToken != nil {
		return nil, errToken
	}

	payload := &AuthPayload{
		Token: &token,
		User:  u,
	}
	return payload, nil
}
