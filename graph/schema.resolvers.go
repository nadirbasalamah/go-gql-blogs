package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/nadirbasalamah/go-gql-blogs/graph/generated"
	"github.com/nadirbasalamah/go-gql-blogs/graph/middleware"
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
)

func (r *mutationResolver) Register(ctx context.Context, input model.NewUser) (string, error) {
	var token string = r.userService.Register(input)

	if token == "" {
		return "", errors.New("registration failed")
	}

	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (string, error) {
	var token string = r.userService.Login(input)

	if token == "" {
		return "", errors.New("login failed, invalid email or password")
	}

	return token, nil
}

func (r *mutationResolver) NewBlog(ctx context.Context, input model.NewBlog) (*model.Blog, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return &model.Blog{}, errors.New("access denied")
	}

	blog, err := r.blogService.CreateBlog(input, *user)

	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (r *mutationResolver) EditBlog(ctx context.Context, input model.EditBlog) (*model.Blog, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return &model.Blog{}, errors.New("access denied")
	}

	blog, err := r.blogService.EditBlog(input, *user)
	if err != nil {
		return &model.Blog{}, err
	}

	return blog, nil
}

func (r *mutationResolver) DeleteBlog(ctx context.Context, input model.DeleteBlog) (bool, error) {
	user := middleware.ForContext(ctx)
	if user == nil {
		return false, errors.New("access denied")
	}

	var result bool = r.blogService.DeleteBlog(input, *user)

	return result, nil
}

func (r *queryResolver) Blogs(ctx context.Context) ([]*model.Blog, error) {
	var blogs []*model.Blog = r.blogService.GetAllBlogs()

	return blogs, nil
}

func (r *queryResolver) Blog(ctx context.Context, id string) (*model.Blog, error) {
	blog, err := r.blogService.GetBlogByID(id)

	if err != nil {
		return &model.Blog{}, err
	}

	return blog, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
