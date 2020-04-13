package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/models"
	"github.com/antonivlev/gql-server/resolvers"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
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
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	errDb := database.SetupDatabase()
	if errDb != nil {
		panic(errDb)
	}

	// initalise root resolver
	res := &resolvers.RootResolver{}
	res.NewLinks = make(chan *models.Link)

	schema := parseSchema("./schema.graphql", res)
	graphQLHandler := graphqlws.NewHandlerFunc(schema, &relay.Handler{Schema: schema})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// pass auth token to resolvers
		token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		ctx := context.WithValue(context.Background(), "token", token)
		graphQLHandler(w, r.WithContext(ctx))
	})

	log.Println("serving on 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
