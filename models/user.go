package models

import (
	"context"
	"errors"
	"fmt"
	"loginapp/utils/token"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	User_ID  string             `json:"user_id" bson:"user_id,omitempty"`
}

func GetUserByID(uid string) (User, error) {

	var u User

	collection := MongoDB.Collection("Users")
	ctx := context.TODO()

	err := collection.FindOne(ctx, User{User_ID: uid}).Decode(&u)
	if err != nil {
		return u, errors.New("User not found!")
	}

	u.PrepareGive()

	return u, nil

}

func (u *User) PrepareGive() {
	u.Password = ""
}

func VerifyPassword(password, hashedPassword string) bool {
	return hashedPassword == password
}

func LoginCheck(username string, password string) (string, error) {

	var err error

	u := User{}
	ctx := context.TODO()
	collection := MongoDB.Collection("Users")
	err = collection.FindOne(ctx, User{Username: username}).Decode(&u)

	if err != nil {
		return "", err
	}

	var isCorrectPassword bool = VerifyPassword(password, u.Password)

	if !isCorrectPassword {
		return "", fmt.Errorf("WrongPassword")
	}

	token, err := token.GenerateToken(u.User_ID)

	if err != nil {
		return "", err
	}

	return token, nil

}

func (u *User) SaveUser() (*User, error) {

	var err error
	collection := MongoDB.Collection("Users")
	ctx := context.TODO()
	collection.InsertOne(ctx, u)
	return u, err
}
