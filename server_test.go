package main

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/nadirbasalamah/go-gql-blogs/database"
	"github.com/nadirbasalamah/go-gql-blogs/graph"
	"github.com/nadirbasalamah/go-gql-blogs/graph/generated"
	"github.com/nadirbasalamah/go-gql-blogs/graph/middleware"
	"github.com/steinfletcher/apitest"
)

func graphQLHandler() *chi.Mux {
	var router *chi.Mux = chi.NewRouter()
	router.Use(middleware.NewMiddleware())

	err := database.Connect()
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v\n", err)
	}

	fmt.Println("Connected to the database")

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	return router
}

func TestGetBlogs_Success(t *testing.T) {
	apitest.New().
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(`query { blogs { title } }`).
		Expect(t).
		Status(http.StatusOK).
		End()
}
