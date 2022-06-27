package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
)

type BlogService struct{}

var storage []*model.Blog

func (b *BlogService) GetAllBlogs() []*model.Blog {
	return storage
}
func (b *BlogService) GetBlogByID(id string) (*model.Blog, error) {
	for _, blog := range storage {
		if blog.ID == id {
			return blog, nil
		}
	}

	return &model.Blog{}, errors.New("blog not found")
}
func (b *BlogService) CreateBlog(input model.NewBlog) *model.Blog {
	var blog model.Blog = model.Blog{
		ID:      uuid.New().String(),
		Title:   input.Title,
		Content: input.Content,
	}

	storage = append(storage, &blog)

	return &blog
}
func (b *BlogService) EditBlog(input model.EditBlog) (*model.Blog, error) {
	for _, blog := range storage {
		if blog.ID == input.BlogID {
			blog.Title = input.Title
			blog.Content = input.Content

			return blog, nil
		}
	}

	return nil, errors.New("blog not found")
}

func (b *BlogService) DeleteBlog(input model.DeleteBlog) bool {
	var afterDeleted []*model.Blog = []*model.Blog{}

	for _, blog := range storage {
		if blog.ID != input.BlogID {
			afterDeleted = append(afterDeleted, blog)
		}
	}

	storage = afterDeleted

	return true
}
