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
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
	"github.com/nadirbasalamah/go-gql-blogs/utils"
	"github.com/steinfletcher/apitest"
)

func graphQLHandler() *chi.Mux {
	var router *chi.Mux = chi.NewRouter()
	router.Use(middleware.NewMiddleware())

	err := database.ConnectTest()
	if err != nil {
		log.Fatalf("Cannot connect to the test database: %v\n", err)
	}

	fmt.Println("Connected to the test database")

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	return router
}

func getJWTToken(user model.User) string {
	token, err := utils.GenerateNewAccessToken(user.ID)
	if err != nil {
		panic(err)
	}

	return "Bearer " + token
}

func getBlog() model.Blog {
	database.ConnectTest()
	blog, err := database.SeedBlog()
	if err != nil {
		panic(err)
	}

	return blog
}

func getUser() model.User {
	database.ConnectTest()
	user, err := database.SeedUser()
	if err != nil {
		panic(err)
	}

	return user
}

func TestSignup_Success(t *testing.T) {
	apitest.New().
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(`mutation {
			register(input:{
				email:"test@test.com",
				username:"test",
				password:"123123"
			})
		}`).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestLogin_Success(t *testing.T) {
	database.ConnectTest()
	user, err := database.SeedUser()
	if err != nil {
		panic(err)
	}

	var query string = `mutation {
		login(input:{
			email:"` + user.Email + `",
			password:"` + user.Password + `"
		})
	}`

	apitest.New().
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		End()
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

func TestGetBlog_Success(t *testing.T) {
	var blog model.Blog = getBlog()

	var query string = `query {
		blog(id:"` + blog.ID + `") {
			title
			content
			createdAt
			updatedAt
		}
	}`

	apitest.New().
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestCreateBlog_Success(t *testing.T) {
	var author model.User = getUser()
	var token string = getJWTToken(author)

	var query string = `mutation {
		newBlog(input:{
			title:"my blog",
			content:"this is the content"
		}) {
			id
			title
		}
	}`

	apitest.New().
		Handler(graphQLHandler()).
		Post("/query").
		Header("Authorization", token).
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		End()
}
