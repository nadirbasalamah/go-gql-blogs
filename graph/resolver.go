package graph

import "github.com/nadirbasalamah/go-gql-blogs/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	blogs []*model.Blog
}
