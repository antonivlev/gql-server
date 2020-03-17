package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type RootResolver struct{}

var links = []models.Link{
	models.Link{
		URL:         "www.bbc.co.uk",
		Description: "I'm a link!",
	},
}

func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil
}

func (r *RootResolver) Feed() ([]models.Link, error) {
	return links, nil
}

func (r *RootResolver) Post(args struct {
	Description string
	URL         string
}) (models.Link, error) {
	newLink := models.Link{
		URL:         args.URL,
		Description: args.Description,
	}
	links = append(links, newLink)

	// this should return the newly created link, with its id
	dbLink, errCreate := database.CreateLink(newLink)
	if errCreate != nil {
		return models.Link{}, errCreate
	}
	return *dbLink, nil
}

var (
	// We can pass an option to the schema so we don’t need to
	// write a method to access each type’s field:
	opts = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
)

// Reads and parses the schema from file. Associates resolver. Panics if can't read.
func parseSchema(path string, resolver interface{}) *graphql.Schema {
	bstr, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	schemaString := string(bstr)
	parsedSchema, err := graphql.ParseSchema(schemaString, resolver, opts...)
	if err != nil {
		panic(err)
	}
	return parsedSchema
}

func main() {
	errDb := database.SetupDatabase()
	if errDb != nil {
		panic(errDb)
	}

	http.Handle("/query", &relay.Handler{
		Schema: parseSchema("./schema.graphql", &RootResolver{}),
	})
	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
