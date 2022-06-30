package service

import (
	"context"
	"errors"

	"github.com/nadirbasalamah/go-gql-blogs/database"
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlogService struct{}

const BLOG_COLLECTION = "blogs"

func (b *BlogService) GetAllBlogs() []*model.Blog {
	var query primitive.D = bson.D{{}}

	cursor, err := database.GetCollection(BLOG_COLLECTION).Find(context.TODO(), query)
	if err != nil {
		return []*model.Blog{}
	}

	var blogs []*model.Blog = make([]*model.Blog, 0)

	if err := cursor.All(context.TODO(), &blogs); err != nil {
		return []*model.Blog{}
	}

	return blogs
}

func (b *BlogService) GetBlogByID(id string) (*model.Blog, error) {
	blogID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &model.Blog{}, errors.New("blog not found")
	}

	var query primitive.D = bson.D{{Key: "_id", Value: blogID}}
	var collection *mongo.Collection = database.GetCollection(BLOG_COLLECTION)

	var blogData *mongo.SingleResult = collection.FindOne(context.TODO(), query)

	var blog *model.Blog = &model.Blog{}
	blogData.Decode(blog)

	return blog, nil
}

func (b *BlogService) CreateBlog(input model.NewBlog) (*model.Blog, error) {
	var blog model.Blog = model.Blog{
		Title:   input.Title,
		Content: input.Content,
	}

	var collection *mongo.Collection = database.GetCollection(BLOG_COLLECTION)

	result, err := collection.InsertOne(context.TODO(), blog)

	if err != nil {
		return &model.Blog{}, errors.New("create blog failed")
	}

	var filter primitive.D = bson.D{{Key: "_id", Value: result.InsertedID}}
	var createdRecord *mongo.SingleResult = collection.FindOne(context.TODO(), filter)

	var createdBlog *model.Blog = &model.Blog{}

	createdRecord.Decode(createdBlog)

	return createdBlog, nil

}

func (b *BlogService) EditBlog(input model.EditBlog) (*model.Blog, error) {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)
	if err != nil {
		return &model.Blog{}, errors.New("id is invalid")
	}

	var query primitive.D = bson.D{{Key: "_id", Value: blogID}}
	var update primitive.D = bson.D{{
		Key: "$set",
		Value: bson.D{
			{Key: "title", Value: input.Title},
			{Key: "content", Value: input.Content},
		},
	}}

	var collection *mongo.Collection = database.GetCollection(BLOG_COLLECTION)

	var updateResult *mongo.SingleResult = collection.FindOneAndUpdate(context.TODO(), query, update)

	if updateResult.Err() != nil {
		if err == mongo.ErrNoDocuments {
			return &model.Blog{}, errors.New("blog not found")
		}
		return &model.Blog{}, errors.New("update blog failed")
	}

	blog, _ := b.GetBlogByID(input.BlogID)
	return blog, nil
}

func (b *BlogService) DeleteBlog(input model.DeleteBlog) bool {
	blogID, err := primitive.ObjectIDFromHex(input.BlogID)
	if err != nil {
		return false
	}

	var query primitive.D = bson.D{{Key: "_id", Value: blogID}}
	var collection *mongo.Collection = database.GetCollection(BLOG_COLLECTION)

	result, err := collection.DeleteOne(context.TODO(), query)
	var isFailed bool = err != nil || result.DeletedCount < 1

	if isFailed {
		return !isFailed
	}

	return true
}
