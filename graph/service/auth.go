package service

import (
	"context"
	"errors"
	"time"

	"github.com/nadirbasalamah/go-gql-blogs/database"
	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
	"github.com/nadirbasalamah/go-gql-blogs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct{}

const USER_COLLECTION = "users"

func (u *UserService) Register(input model.NewUser) string {
	bs, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}

	var password string = string(bs)

	var user model.User = model.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	var collection *mongo.Collection = database.GetCollection(USER_COLLECTION)

	res, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		return ""
	}

	var userId string = res.InsertedID.(primitive.ObjectID).Hex()

	token, err := utils.GenerateNewAccessToken(userId)

	if err != nil {
		return ""
	}

	return token
}

func (u *UserService) Login(input model.LoginInput) string {
	var collection *mongo.Collection = database.GetCollection(USER_COLLECTION)

	var user *model.User = &model.User{}
	filter := bson.M{"email": input.Email}

	var res *mongo.SingleResult = collection.FindOne(context.TODO(), filter)
	if err := res.Decode(user); err != nil {
		return ""
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return ""
	}

	token, err := utils.GenerateNewAccessToken(user.ID)

	if err != nil {
		return ""
	}

	return token
}

func (u *UserService) GetUser(id string) (*model.User, error) {
	userID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &model.User{}, errors.New("id is invalid")
	}

	var query primitive.D = bson.D{{Key: "_id", Value: userID}}
	var collection *mongo.Collection = database.GetCollection(USER_COLLECTION)

	var userData *mongo.SingleResult = collection.FindOne(context.TODO(), query)

	if userData.Err() != nil {
		return &model.User{}, errors.New("user not found")
	}

	var user *model.User = &model.User{}
	userData.Decode(user)

	return user, nil
}
