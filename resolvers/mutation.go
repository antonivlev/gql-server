package resolvers

import (
	"context"
	"fmt"

	"github.com/antonivlev/gql-server/auth"
	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
	"github.com/graph-gophers/graphql-go"
)

func (r *RootResolver) Post(ctx context.Context, args struct {
	Description string
	URL         string
}) (models.Link, error) {
	poster, errUser := database.GetUser(ctx)
	if errUser != nil {
		return models.Link{}, errUser
	}

	newLink := models.Link{
		URL:         args.URL,
		Description: args.Description,
		PostedBy:    poster,
	}

	dbLink, errCreate := database.CreateLink(newLink)
	if errCreate != nil {
		return models.Link{}, errCreate
	}

	select {
	case r.NewLinks <- dbLink:
		// values are being read from r.Events
		fmt.Println("r.NewLinks: inserted link")
	default:
		// no subscribers, link not in channel
		fmt.Println("r.Links: link created, not inserted")
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

	token, errToken := auth.GenerateToken(u)
	if errToken != nil {
		return nil, errToken
	}

	payload := &AuthPayload{
		Token: &token,
		User:  u,
	}
	return payload, nil
}

func (r *RootResolver) Vote(
	ctx context.Context,
	args struct{ LinkID graphql.ID },
) (*models.Vote, error) {
	voter, errUser := database.GetUser(ctx)
	if errUser != nil {
		return nil, errUser
	}

	dbVote, errCreate := database.CreateVote(voter.ID, args.LinkID)
	if errCreate != nil {
		return nil, errCreate
	}

	select {
	case r.NewVotes <- dbVote:
		fmt.Println("r.NewVotes: inserted vote")
	default:
		fmt.Println("r.NewVotes: vote created, not inserted")
	}

	return dbVote, nil
}
