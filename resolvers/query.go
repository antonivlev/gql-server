package resolvers

import (
	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
)

func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil
}

type Feed struct {
	Links []models.Link
	Count int32
}

func (r *RootResolver) Feed(
	args struct {
		Filter  *string
		Skip    *int32
		First   *int32
		OrderBy *string
	}) (Feed, error) {

	links, count, err := database.GetAllLinks(args.Filter, args.OrderBy, args.Skip, args.First)
	if err != nil {
		return Feed{}, err
	}

	feed := Feed{
		Links: links,
		Count: int32(count),
	}

	return feed, nil
}
