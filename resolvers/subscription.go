package resolvers

import (
	"fmt"

	"github.com/antonivlev/gql-server/models"
)

func (r *RootResolver) NewLink() (chan *models.Link, error) {
	fmt.Println("subcribing to links")
	return r.NewLinks, nil
}

func (r *RootResolver) NewVote() (chan *models.Vote, error) {
	fmt.Println("subcribing to votes")
	return r.NewVotes, nil
}
