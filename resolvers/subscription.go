package resolvers

import (
	"fmt"

	"github.com/antonivlev/gql-server/models"
)

func (r *RootResolver) NewLink() (chan *models.Link, error) {
	fmt.Println("subcribing")
	return r.NewLinks, nil
}
