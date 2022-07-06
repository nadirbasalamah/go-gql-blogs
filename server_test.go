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

func cleanup(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
	if http.StatusOK == res.StatusCode {
		database.CleanSeeders()
	}
}

func TestSignup_Success(t *testing.T) {
	apitest.New().
		Observe(cleanup).
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
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		End()
}

func TestLogin_Failed(t *testing.T) {
	var result string = `{
		"errors": [
			{
				"message": "login failed, invalid email or password",
				"path": [
					"login"
				]
			}
		],
		"data": null
	}`

	var query string = `mutation {
		login(input:{
			email:"wrong@mail.com",
			password:"123456"
		})
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
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
		}
	}`

	var result string = `{
		"data": {
			"blog": {
				"title": "` + blog.Title + `",
				"content": "` + blog.Content + `"
			}
		}
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestGetBlog_Failed(t *testing.T) {
	var query string = `query {
		blog(id:"62c24f3b4896bb25c21e49b9") {
			title
			content
		}
	}`

	var result string = `{
		"errors": [
			{
				"message": "blog not found",
				"path": [
					"blog"
				]
			}
		],
		"data": null
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
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
			title
			content
		}
	}`

	var result string = `{
		"data": {
			"newBlog": {
				"title": "my blog",
				"content": "this is the content"
			}
		}
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		Header("Authorization", token).
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestCreateBlog_Failed(t *testing.T) {
	var query string = `mutation {
		newBlog(input:{
			title:"my blog",
			content:"this is the content"
		}) {
			title
			content
		}
	}`

	var result string = `{
		"errors": [
			{
				"message": "access denied",
				"path": [
					"newBlog"
				]
			}
		],
		"data": null
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestEditBlog_Success(t *testing.T) {
	var blog model.Blog = getBlog()
	var token string = getJWTToken(*blog.Author)

	var query string = `mutation {
		editBlog(input:{
			blogId:"` + blog.ID + `"
			title:"my blog",
			content:"this is the content"
		}) {
			title
			content
		}
	}`

	var result string = `{
		"data": {
			"editBlog": {
				"title": "my blog",
				"content": "this is the content"
			}
		}
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		Header("Authorization", token).
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestEditBlog_Failed(t *testing.T) {
	var blog model.Blog = getBlog()

	var query string = `mutation {
		editBlog(input:{
			blogId:"` + blog.ID + `"
			title:"my blog",
			content:"this is the content"
		}) {
			title
			content
		}
	}`

	var result string = `{
		"errors": [
			{
				"message": "access denied",
				"path": [
					"editBlog"
				]
			}
		],
		"data": null
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestDeleteBlog_Success(t *testing.T) {
	var blog model.Blog = getBlog()
	var token string = getJWTToken(*blog.Author)

	var query string = `mutation {
		deleteBlog(input:{
			blogId:"` + blog.ID + `"
		})
	}`

	var result string = `{
		"data": {
			"deleteBlog": true
		}
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		Header("Authorization", token).
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}

func TestDeleteBlog_Failed(t *testing.T) {
	var blog model.Blog = getBlog()

	var query string = `mutation {
		deleteBlog(input:{
			blogId:"` + blog.ID + `"
		})
	}`

	var result string = `{
		"errors": [
			{
				"message": "access denied",
				"path": [
					"deleteBlog"
				]
			}
		],
		"data": null
	}`

	apitest.New().
		Observe(cleanup).
		Handler(graphQLHandler()).
		Post("/query").
		GraphQLQuery(query).
		Expect(t).
		Status(http.StatusOK).
		Body(result).
		End()
}
