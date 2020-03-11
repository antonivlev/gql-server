package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

type RootResolver struct{}

type Link struct {
	ID          graphql.ID
	Description string
	URL         string
}

var links = []Link{
	Link{
		ID:          "link-0",
		URL:         "www.bbc.co.uk",
		Description: "I'm a link!",
	},
}

func (r *RootResolver) Info() (string, error) {
	return "this is a thing", nil
}

func (r *RootResolver) Feed() ([]Link, error) {
	return links, nil
}

func (r *RootResolver) Post(args struct {
	Description string
	URL         string
}) (Link, error) {
	newLink := Link{
		ID:          graphql.ID("link-" + fmt.Sprint(len(links))),
		URL:         args.URL,
		Description: args.Description,
	}
	links = append(links, newLink)
	return newLink, nil
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
	http.Handle("/query", &relay.Handler{
		Schema: parseSchema("./schema.graphql", &RootResolver{}),
	})
	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
