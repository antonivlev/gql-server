package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/resolvers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

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
		Schema: parseSchema("./schema.graphql", &resolvers.RootResolver{}),
	})
	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
