package database

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/nadirbasalamah/go-gql-blogs/graph/model"
	"github.com/nadirbasalamah/go-gql-blogs/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type MongoInstance struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var DB MongoInstance

func Connect() error {
	client, _ := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	var db *mongo.Database = client.Database(os.Getenv("DATABASE_NAME"))

	if err != nil {
		return err
	}

	DB = MongoInstance{
		Client:   client,
		Database: db,
	}

	return nil
}

func ConnectTest() error {
	client, _ := mongo.NewClient(options.Client().ApplyURI(utils.GetValue("MONGO_URI")))
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := client.Connect(ctx)
	var db *mongo.Database = client.Database(utils.GetValue("DATABASE_TEST_NAME"))

	if err != nil {
		return err
	}

	DB = MongoInstance{
		Client:   client,
		Database: db,
	}

	return nil
}

func GetCollection(name string) *mongo.Collection {
	return DB.Database.Collection(name)
}

func SeedBlog() (model.Blog, error) {
	blogFaker, err := utils.CreateFaker[utils.BlogFaker]()
	if err != nil {
		return model.Blog{}, err
	}

	author, err := SeedUser()
	if err != nil {
		return model.Blog{}, err
	}

	var blog model.Blog = model.Blog{
		Title:     blogFaker.Title,
		Content:   blogFaker.Content,
		Author:    &author,
		CreatedAt: time.Now(),
	}

	var collection *mongo.Collection = GetCollection("blogs")

	result, err := collection.InsertOne(context.TODO(), blog)

	if err != nil {
		return model.Blog{}, errors.New("create blog failed")
	}

	var filter primitive.D = bson.D{{Key: "_id", Value: result.InsertedID}}
	var createdRecord *mongo.SingleResult = collection.FindOne(context.TODO(), filter)

	var createdBlog *model.Blog = &model.Blog{}

	createdRecord.Decode(createdBlog)

	return *createdBlog, nil
}

func SeedUser() (model.User, error) {
	userFaker, err := utils.CreateFaker[utils.UserFaker]()
	if err != nil {
		return model.User{}, err
	}

	bs, err := bcrypt.GenerateFromPassword([]byte(userFaker.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	var password string = string(bs)

	var user model.User = model.User{
		Username:  userFaker.Username,
		Email:     userFaker.Email,
		Password:  password,
		CreatedAt: time.Now(),
	}

	var collection *mongo.Collection = GetCollection("users")

	result, err := collection.InsertOne(context.TODO(), user)

	if err != nil {
		return model.User{}, err
	}

	user.Password = userFaker.Password
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}

func CleanSeeders() {
	var userCollection *mongo.Collection = GetCollection("users")
	var blogCollection *mongo.Collection = GetCollection("blogs")

	userErr := userCollection.Drop(context.TODO())
	blogErr := blogCollection.Drop(context.TODO())

	var isFailed bool = userErr != nil || blogErr != nil

	if isFailed {
		panic("error when deleting all data inside collection")
	}
}
