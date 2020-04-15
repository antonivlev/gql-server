package resolvers

import "github.com/antonivlev/gql-server/models"

type RootResolver struct {
	NewLinks chan *models.Link
	NewVotes chan *models.Vote
}
