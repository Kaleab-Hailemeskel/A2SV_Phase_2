package data

import (
	"context"
	"fmt"
	"task-6_authentication_and_authorization/models"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	userDataBaseName   = "UserBase"
	userCollectionName = "users"
)

func FindOneUser(userEmail string) (*models.User, error) {
	collection := mongoClient.Database(userDataBaseName).Collection(userCollectionName)
	filter := bson.M{"email": userEmail}
	var user models.User
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func InsertOneUser(user models.User) error {
	collection := mongoClient.Database(userDataBaseName).Collection(userCollectionName)
	inserted, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	fmt.Println("User Registered with ID", inserted.InsertedID)
	return nil
}
