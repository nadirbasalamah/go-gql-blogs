package utils

import (
	"time"

	"github.com/bxcodec/faker/v3"
)

type User struct {
	ID        string     `json:"id" bson:"_id,omitempty"`
	Username  string     `json:"username" bson:"username" faker:"username"`
	Email     string     `json:"email" bson:"email" faker:"email"`
	Password  string     `json:"password" bson:"password" faker:"password"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Blog struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Title     string    `json:"title" bson:"title" faker:"name"`
	Content   string    `json:"content" bson:"content" faker:"word"`
	Author    *User     `json:"author" bson:"author"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func CreateFaker[T any]() (T, error) {
	var fakerData *T = new(T)
	err := faker.FakeData(fakerData)
	if err != nil {
		return *fakerData, err
	}

	return *fakerData, nil
}
