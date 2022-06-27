package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nadirbasalamah/go-gql-blogs/graph/generated"
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
)

func (r *mutationResolver) NewBlog(ctx context.Context, input model.NewBlog) (*model.Blog, error) {
	var blog model.Blog = model.Blog{
		ID:      uuid.New().String(),
		Title:   input.Title,
		Content: input.Content,
	}

	r.blogs = append(r.blogs, &blog)

	return &blog, nil
}

func (r *mutationResolver) EditBlog(ctx context.Context, input model.EditBlog) (*model.Blog, error) {
	for _, blog := range r.blogs {
		if blog.ID == input.BlogID {
			blog.Title = input.Title
			blog.Content = input.Content

			return blog, nil
		}
	}

	return nil, errors.New("blog not found")
}

func (r *mutationResolver) DeleteBlog(ctx context.Context, input model.DeleteBlog) (bool, error) {
	var afterDeleted []*model.Blog = []*model.Blog{}

	for _, blog := range r.blogs {
		if blog.ID != input.BlogID {
			afterDeleted = append(afterDeleted, blog)
		}
	}

	r.blogs = afterDeleted

	return true, nil
}

func (r *queryResolver) Blogs(ctx context.Context) ([]*model.Blog, error) {
	return r.blogs, nil
}

func (r *queryResolver) Blog(ctx context.Context, id string) (*model.Blog, error) {
	for _, blog := range r.blogs {
		if blog.ID == id {
			return blog, nil
		}
	}

	return nil, errors.New("blog not found")
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
// func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
// 	panic(fmt.Errorf("not implemented"))
// }
// func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
// 	panic(fmt.Errorf("not implemented"))
// }
