package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/antonivlev/gql-server/database"
	"github.com/antonivlev/gql-server/resolvers"
	graphql "github.com/graph-gophers/graphql-go"
	"golang.org/x/net/websocket"
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

	var conn *websocket.Conn

	type message struct {
		Message string `json:"message"`
	}
	socket := func(ws *websocket.Conn) {
		conn = ws
		for {
			// allocate our container struct
			var m message
			// receive a message using the codec
			if err := websocket.JSON.Receive(ws, &m); err != nil {
				log.Println(err)
			}
			log.Println("Received message:", m.Message)
			// send a response
			m2 := message{"Thanks for the message!"}
			if err := websocket.JSON.Send(ws, m2); err != nil {
				log.Println(err)
			}
		}
	}
	http.Handle("/socket", websocket.Handler(socket))

	schema := parseSchema("./schema.graphql", &resolvers.RootResolver{})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := websocket.JSON.Send(conn, message{"oi op"}); err != nil {
			log.Println(err)
		}

		type Payload struct {
			Query         string
			OperationName string
			Variables     map[string]interface{}
		}
		var payload Payload

		errParse := json.NewDecoder(r.Body).Decode(&payload)
		if errParse != nil {
			http.Error(w, errParse.Error(), http.StatusBadRequest)
			return
		}

		token := strings.ReplaceAll(r.Header.Get("Authorization"), "Bearer ", "")
		ctx := context.WithValue(context.Background(), "token", token)

		resp := schema.Exec(ctx, payload.Query, payload.OperationName, payload.Variables)

		if len(resp.Errors) > 0 {
			type ErrorResponse struct {
				Data   interface{} `json:"data"`
				Errors interface{} `json:"errors"`
			}

			errResp := ErrorResponse{
				Data:   resp.Data,
				Errors: resp.Errors,
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errResp)
			return
		}
		json, err := json.MarshalIndent(resp, "", "\t")
		if err != nil {
			log.Printf("json.MarshalIndent: %s", err)
			return
		}

		fmt.Fprint(w, string(json))
	})

	fmt.Println("serving on 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
